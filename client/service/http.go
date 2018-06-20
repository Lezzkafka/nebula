package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/grpc/status"

	"time"

	"github.com/rs/cors"
	"github.com/samoslab/nebula/client/common"
	"github.com/samoslab/nebula/client/config"
	"github.com/samoslab/nebula/client/daemon"
	regclient "github.com/samoslab/nebula/client/register"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"golang.org/x/crypto/acme/autocert"
)

const (
	shutdownTimeout = time.Second * 5

	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// The timeout configuration is necessary for public servers, or else
	// connections will be used up
	serverReadTimeout  = time.Second * 10
	serverWriteTimeout = time.Second * 60
	serverIdleTimeout  = time.Second * 120

	// Directory where cached SSL certs from Let's Encrypt are stored
	tlsAutoCertCache = "cert-cache"
)

var (
	errInternalServerError = errors.New("Internal Server Error")
)

// HTTPServer exposes the API endpoints and static website
type HTTPServer struct {
	cfg config.Config
	log logrus.FieldLogger
	//log           *logrus.Logger
	cm            *daemon.ClientManager
	httpListener  *http.Server
	httpsListener *http.Server
	quit          chan struct{}
	done          chan struct{}
}

func InitClientManager(log logrus.FieldLogger, webcfg config.Config) (*daemon.ClientManager, error) {
	_, defaultConfig := daemon.GetConfigFile()
	clientConfig, err := config.LoadConfig(defaultConfig)
	if err != nil {
		if err == config.ErrNoConf {
			return nil, fmt.Errorf("Config file is not ready, please call /store/register to register first")
		} else if err == config.ErrConfVerify {
			return nil, fmt.Errorf("Config file wrong, can not start daemon.")
		}
	}
	cm, err := daemon.NewClientManager(log, webcfg, clientConfig)
	if err != nil {
		fmt.Printf("new client manager failed %v\n", err)
		return cm, err
	}
	return cm, nil
}

// NewHTTPServer creates an HTTPServer
func NewHTTPServer(log logrus.FieldLogger, cfg config.Config) *HTTPServer {
	cm, err := InitClientManager(log, cfg)
	if err != nil {
		log.Errorf("init client manager failed, error %v", err)
	}
	return &HTTPServer{
		cfg:  cfg,
		log:  log,
		cm:   cm,
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
}

func (s *HTTPServer) CanBeWork() bool {
	if s.cm != nil {
		return true
	}
	return false
}

// Run runs the HTTPServer
func (s *HTTPServer) Run() error {
	log := s.log
	log.Info("HTTP service start")
	defer log.Info("HTTP service closed")
	defer close(s.done)

	var mux http.Handler = s.setupMux()

	allowedHosts := []string{} // empty array means all hosts allowed
	sslHost := ""
	if s.cfg.AutoTLSHost == "" {
		// Note: if AutoTLSHost is not set, but HTTPSAddr is set, then
		// http will redirect to the HTTPSAddr listening IP, which would be
		// either 127.0.0.1 or 0.0.0.0
		// When running behind a DNS name, make sure to set AutoTLSHost
		sslHost = s.cfg.HTTPSAddr
	} else {
		sslHost = s.cfg.AutoTLSHost
		// When using -auto-tls-host,
		// which implies automatic Let's Encrypt SSL cert generation in production,
		// restrict allowed hosts to that host.
		allowedHosts = []string{s.cfg.AutoTLSHost}
	}

	if len(allowedHosts) == 0 {
		log = log.WithField("allowedHosts", "all")
	} else {
		log = log.WithField("allowedHosts", allowedHosts)
	}

	log = log.WithField("sslHost", sslHost)

	log.Info("Configured")

	secureMiddleware := configureSecureMiddleware(sslHost, allowedHosts)
	mux = secureMiddleware.Handler(mux)

	if s.cfg.HTTPAddr != "" {
		s.httpListener = setupHTTPListener(s.cfg.HTTPAddr, mux)
	}

	handleListenErr := func(f func() error) error {
		if err := f(); err != nil {
			select {
			case <-s.quit:
				return nil
			default:
				log.WithError(err).Error("ListenAndServe or ListenAndServeTLS error")
				return fmt.Errorf("http serve failed: %v", err)
			}
		}
		return nil
	}

	if s.cfg.HTTPAddr != "" {
		log.Info(fmt.Sprintf("HTTP server listening on http://%s", s.cfg.HTTPAddr))
	}
	if s.cfg.HTTPSAddr != "" {
		log.Info(fmt.Sprintf("HTTPS server listening on https://%s", s.cfg.HTTPSAddr))
	}

	var tlsCert, tlsKey string
	if s.cfg.HTTPSAddr != "" {
		log.Info("Using TLS")

		s.httpsListener = setupHTTPListener(s.cfg.HTTPSAddr, mux)

		tlsCert = s.cfg.TLSCert
		tlsKey = s.cfg.TLSKey

		if s.cfg.AutoTLSHost != "" {
			log.Info("Using Let's Encrypt autocert")
			// https://godoc.org/golang.org/x/crypto/acme/autocert
			// https://stackoverflow.com/a/40494806
			certManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(s.cfg.AutoTLSHost),
				Cache:      autocert.DirCache(tlsAutoCertCache),
			}

			s.httpsListener.TLSConfig = &tls.Config{
				GetCertificate: certManager.GetCertificate,
			}

			// These will be autogenerated by the autocert middleware
			tlsCert = ""
			tlsKey = ""
		}

	}

	return handleListenErr(func() error {
		var wg sync.WaitGroup
		errC := make(chan error)

		if s.cfg.HTTPAddr != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := s.httpListener.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.WithError(err).Println("ListenAndServe error")
					errC <- err
				}
			}()
		}

		if s.cfg.HTTPSAddr != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := s.httpsListener.ListenAndServeTLS(tlsCert, tlsKey); err != nil && err != http.ErrServerClosed {
					log.WithError(err).Error("ListenAndServeTLS error")
					errC <- err
				}
			}()
		}

		done := make(chan struct{})

		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case err := <-errC:
			return err
		case <-s.quit:
			return nil
		case <-done:
			return nil
		}
	})
}

func configureSecureMiddleware(sslHost string, allowedHosts []string) *secure.Secure {
	sslRedirect := true
	if sslHost == "" {
		sslRedirect = false
	}

	return secure.New(secure.Options{
		AllowedHosts: allowedHosts,
		SSLRedirect:  sslRedirect,
		SSLHost:      sslHost,

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
		// FIXME: Web frontend code has inline styles, CSP doesn't work yet
		// ContentSecurityPolicy: "default-src 'self'",

		// Set HSTS to one year, for this domain only, do not add to chrome preload list
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
		STSSeconds:           31536000, // 1 year
		STSIncludeSubdomains: false,
		STSPreload:           false,

		// Deny use in iframes
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
		FrameDeny: true,

		// Disable MIME sniffing in browsers
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options
		ContentTypeNosniff: true,

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
		BrowserXssFilter: true,

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy
		// "same-origin" is invalid in chrome
		ReferrerPolicy: "no-referrer",
	})
}

func setupHTTPListener(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}
}

func (s *HTTPServer) setupMux() *http.ServeMux {
	mux := http.NewServeMux()

	handleAPI := func(path string, h http.Handler) {
		// Allow requests from a local samos client
		h = cors.New(cors.Options{
			AllowedOrigins: []string{"http://127.0.0.1:7788"},
		}).Handler(h)

		//h = gziphandler.GzipHandler(h)

		mux.Handle(path, h)
	}

	// API Methods
	handleAPI("/api/v1/store/register", RegisterHandler(s))
	handleAPI("/api/v1/store/verifyemail", EmailHandler(s))
	handleAPI("/api/v1/store/folder/add", MkfolderHandler(s))
	handleAPI("/api/v1/store/upload", UploadHandler(s))
	handleAPI("/api/v1/store/download", DownloadHandler(s))
	handleAPI("/api/v1/store/list", ListHandler(s))
	handleAPI("/api/v1/store/remove", RemoveHandler(s))
	handleAPI("/api/v1/store/progress", ProgressHandler(s))
	handleAPI("/api/v1/store/uploaddir", UploadDirHandler(s))
	handleAPI("/api/v1/store/downloaddir", DownloadDirHandler(s))

	handleAPI("/api/v1/package/all", GetAllPackageHandler(s))
	handleAPI("/api/v1/package", GetPackageInfoHandler(s))
	handleAPI("/api/v1/package/buy", BuyPackageHandler(s))
	handleAPI("/api/v1/package/discount", DiscountPackageHandler(s))
	handleAPI("/api/v1/order/all", MyAllOrderHandler(s))
	handleAPI("/api/v1/order/getinfo", GetOrderInfoHandler(s))
	handleAPI("/api/v1/order/recharge/address", RechargeAddressHandler(s))
	handleAPI("/api/v1/order/pay", PayOrderHandler(s))
	handleAPI("/api/v1/usage/amount", UsageAmountHandler(s))

	handleAPI("/api/v1/secret/encrypt", EncryFileHandler(s))
	handleAPI("/api/v1/secret/decrypt", DecryFileHandler(s))

	return mux
}

type RegisterReq struct {
	Email  string `json:"email"`
	Resend bool   `josn:"resend"`
}

type VerifyEmailReq struct {
	Code string `json:"code"`
}

type MkfolderReq struct {
	Folders     []string `json:"folders"`
	Parent      string   `json:"parent"`
	Interactive bool     `json:"interactive"`
}

type UploadReq struct {
	Filename    string `json:"filename"`
	Interactive bool   `json:"interactive"`
	NewVersion  bool   `json:"newversion"`
}

type UploadDirReq struct {
	Parent      string `json:"parent"`
	Interactive bool   `json:"interactive"`
	NewVersion  bool   `json:"newversion"`
}

type DownloadDirReq struct {
	Parent string `json:"parent"`
}

type DownloadReq struct {
	FileHash string `json:"filehash"`
	FileSize uint64 `json:"filesize"`
	FileName string `json:"filename"`
}

type ListReq struct {
	Path     string `json:"path"`
	PageSize uint32 `json:"pagesize"`
	PageNum  uint32 `json:"pagenum"`
	SortType string `json:"sorttype"`
	AscOrder bool   `json:"ascorder"`
}

type RemoveReq struct {
	Target    string `json:"target"`
	Recursion bool   `json:"recursion"`
	IsPath    bool   `json:"ispath"`
}

type ProgressReq struct {
	Files []string `json:"files"`
}

type ProgressRsp struct {
	Progress map[string]float64 `json:"progress"`
}

type EncryFileReq struct {
	FileName   string `json:"file"`
	Password   string `json:"password"`
	OutputFile string `json:output_file`
}

type DecryFileReq struct {
	FileName   string `json:"file"`
	Password   string `json:"password"`
	OutputFile string `json:output_file`
}

func RegisterHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := s.log
		ctx := r.Context()
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		regReq := &RegisterReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&regReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()
		if regReq.Email == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument email must not empty"))
			return
		}

		var err error
		if regReq.Resend {
			err = regclient.ResendVerifyCode(s.cfg.ConfigDir, s.cfg.TrackerServer)
		} else {
			log.Infof("register email %s dir %s", regReq.Email, s.cfg.ConfigDir)
			err = regclient.RegisterClient(log, s.cfg.ConfigDir, s.cfg.TrackerServer, regReq.Email)
		}
		code := 0
		errmsg := ""
		result := "ok"
		if err != nil {
			log.Errorf("send email %+v error %v", regReq, err)
			code = 1
			errmsg = err.Error()
			result = ""
		}
		if !regReq.Resend {
			cm, err := InitClientManager(log, s.cfg)
			if err != nil {
				code = 1
				errmsg = err.Error()
			} else {
				s.cm = cm
			}
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func EmailHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		mailReq := &VerifyEmailReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&mailReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()

		if mailReq.Code == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument code must not empty"))
			return
		}

		err := regclient.VerifyEmail(s.cfg.ConfigDir, s.cfg.TrackerServer, mailReq.Code)
		code := 0
		errmsg := ""
		result := "ok"
		if err != nil {
			log.Errorf("verify email %+v error %v", mailReq, err)
			code = 1
			errmsg = err.Error()
			result = ""
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

// MkfolderHandler create folders
// Method: POST
// Accept: application/json
// URI: /store/folder/add
// Args:
func MkfolderHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		mkReq := &MkfolderReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&mkReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()
		if mkReq.Parent == "" || len(mkReq.Folders) == 0 {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument parent or folders must not empty"))
			return
		}
		for _, folder := range mkReq.Folders {
			if strings.Contains(folder, "/") {
				errorResponse(ctx, w, http.StatusBadRequest, fmt.Errorf("folder %s contains /", folder))
				return
			}
		}

		log.Infof("mkfolder parent %s folders %+v", mkReq.Parent, mkReq.Folders)
		result, err := s.cm.MkFolder(mkReq.Parent, mkReq.Folders, mkReq.Interactive)
		if err != nil {
			log.Errorf("create folder %+v error %v", mkReq.Parent, mkReq.Folders, err)
		}

		code := 0
		errmsg := ""
		if !result {
			code = 1
			errmsg = err.Error()
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func UploadHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		upReq := &UploadReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&upReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()

		if upReq.Filename == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument filename must not empty"))
			return
		}

		log.Infof("upload files %+v", upReq.Filename)
		err := s.cm.UploadFile(upReq.Filename, upReq.Interactive, upReq.NewVersion)
		st, ok := status.FromError(err)
		if !ok {
			log.Infof("err code %d msg %s", st.Code(), st.Message())
		}
		code := 0
		errmsg := ""
		result := "success"
		if err != nil {
			log.Errorf("upload %+v error %v", upReq, err)
			code = 1
			errmsg = err.Error()
			result = ""
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func UploadDirHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		upReq := &UploadDirReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&upReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()

		if upReq.Parent == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument parent must not empty"))
			return
		}

		log.Infof("upload parent %s", upReq.Parent)
		err := s.cm.UploadDir(upReq.Parent, upReq.Interactive, upReq.NewVersion)
		st, ok := status.FromError(err)
		if !ok {
			log.Infof("err code %d msg %s", st.Code(), st.Message())
		}

		code := 0
		errmsg := ""
		result := "success"
		if err != nil {
			log.Errorf("upload %+v error %v", upReq, err)
			code = 1
			errmsg = err.Error()
			result = ""
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}
func DownloadHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		downReq := &DownloadReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&downReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()

		if downReq.FileHash == "" || downReq.FileSize == 0 || downReq.FileName == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument filehash filesize or filename must not empty"))
			return
		}

		log.Infof("download  %+v", downReq)
		err := s.cm.DownloadFile(downReq.FileName, downReq.FileHash, downReq.FileSize)
		code := 0
		errmsg := ""
		result := "success"
		if err != nil {
			log.Errorf("download files %+v error %v", downReq, err)
			code = 1
			errmsg = err.Error()
			result = ""
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func DownloadDirHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		downReq := &DownloadDirReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&downReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()

		if downReq.Parent == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument parent must not empty"))
			return
		}

		log.Infof("download  %+v", downReq)
		err := s.cm.DownloadDir(downReq.Parent)
		code := 0
		errmsg := ""
		result := "success"
		if err != nil {
			log.Errorf("download files %+v error %v", downReq, err)
			code = 1
			errmsg = err.Error()
			result = ""
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func ListHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		listReq := &ListReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&listReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()
		if listReq.Path == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument path must not empty"))
			return
		}

		log.Infof("list %+v", listReq)
		result, err := s.cm.ListFiles(listReq.Path, listReq.PageSize, listReq.PageNum, listReq.SortType, listReq.AscOrder)
		code := 0
		errmsg := ""
		if err != nil {
			log.Errorf("list files %+v error %v", listReq, err)
			code = 1
			errmsg = err.Error()
			result = nil
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func RemoveHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		rmReq := &RemoveReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&rmReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()
		if rmReq.Target == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument target must not empty"))
			return
		}

		log.Infof("remove %+v", rmReq)
		err := s.cm.RemoveFile(rmReq.Target, rmReq.Recursion, false)
		code := 0
		errmsg := ""
		result := true
		if err != nil {
			log.Errorf("remove files %+v error %v", rmReq, err)
			code = 1
			errmsg = err.Error()
			result = false
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func ProgressHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		progressReq := &ProgressReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&progressReq); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()

		log.Infof("progress %+v", progressReq)
		progressRsp, err := s.cm.GetProgress(progressReq.Files)
		code := 0
		errmsg := ""
		if err != nil {
			log.Errorf("remove files %+v error %v", progressReq, err)
			code = 1
			errmsg = err.Error()
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, progressRsp, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func EncryFileHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		req := &EncryFileReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()
		if req.FileName == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument file must not empty"))
			return
		}
		if req.Password == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument password must not empty"))
			return
		}
		if len(req.Password) != 16 {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("password lenght must 16"))
			return
		}

		log.Infof("encrypt file %+v\n", req.FileName)
		if req.OutputFile == "" {
			req.OutputFile = req.FileName
		}
		err := common.EncryptFile(req.FileName, []byte(req.Password), req.OutputFile)
		code := 0
		errmsg := ""
		result := true
		if err != nil {
			log.Errorf("encrypt file %+v error %v", req.FileName, err)
			code = 1
			errmsg = err.Error()
			result = false
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

func DecryFileHandler(s *HTTPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.CanBeWork() {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("register first"))
			return
		}
		log := s.cm.Log
		w.Header().Set("Accept", "application/json")

		if !validMethod(ctx, w, r, []string{http.MethodPost}) {
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			errorResponse(ctx, w, http.StatusUnsupportedMediaType, errors.New("Invalid content type"))
			return
		}

		req := &DecryFileReq{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			err = fmt.Errorf("Invalid json request body: %v", err)
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}

		defer r.Body.Close()
		if req.FileName == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument file must not empty"))
			return
		}
		if req.Password == "" {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("argument password must not empty"))
			return
		}
		if len(req.Password) != 16 {
			errorResponse(ctx, w, http.StatusBadRequest, errors.New("password lenght must 16"))
			return
		}

		log.Infof("encrypt file %+v\n", req.FileName)
		if req.OutputFile == "" {
			req.OutputFile = req.FileName
		}
		err := common.DecryptFile(req.FileName, []byte(req.Password), req.OutputFile)
		code := 0
		errmsg := ""
		result := true
		if err != nil {
			log.Errorf("decrypt file %+v error %v", req.FileName, err)
			code = 1
			errmsg = err.Error()
			result = false
		}

		rsp, err := common.MakeUnifiedHTTPResponse(code, result, errmsg)
		if err != nil {
			errorResponse(ctx, w, http.StatusBadRequest, err)
			return
		}
		if err := JSONResponse(w, rsp); err != nil {
			fmt.Printf("error %v\n", err)
		}
	}
}

// JSONResponse marshal data into json and write response
func JSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	d, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	_, err = w.Write(d)
	return err
}

func errorResponse(ctx context.Context, w http.ResponseWriter, code int, err error) {
	unifiedres, err := common.MakeUnifiedHTTPResponse(code, "", err.Error())
	if err != nil {
		return
	}
	if err := JSONResponse(w, unifiedres); err != nil {
		fmt.Printf("error response failed")
	}
}

// Shutdown stops the HTTPServer
func (s *HTTPServer) Shutdown() {
	s.log.Info("Shutting down HTTP server(s)")
	defer s.log.Info("Shutdown HTTP server(s)")
	close(s.quit)

	var wg sync.WaitGroup
	wg.Add(2)

	shutdown := func(proto string, ln *http.Server) {
		defer wg.Done()
		if ln == nil {
			return
		}
		log := s.log.WithFields(logrus.Fields{
			"proto":   proto,
			"timeout": shutdownTimeout,
		})

		defer log.Info("Shutdown server")
		log.Info("Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := ln.Shutdown(ctx); err != nil {
			log.WithError(err).Error("HTTP server shutdown error")
		}
	}

	shutdown("HTTP", s.httpListener)
	shutdown("HTTPS", s.httpsListener)

	if s.cm != nil {
		s.cm.Shutdown()
	}
	wg.Wait()

	<-s.done
}

func validMethod(ctx context.Context, w http.ResponseWriter, r *http.Request, allowed []string) bool {
	for _, m := range allowed {
		if r.Method == m {
			return true
		}
	}

	w.Header().Set("Allow", strings.Join(allowed, ", "))

	status := http.StatusMethodNotAllowed
	errorResponse(ctx, w, status, errors.New("Invalid request method"))

	return false
}
