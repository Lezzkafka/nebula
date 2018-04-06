// Code generated by protoc-gen-go. DO NOT EDIT.
// source: provider_register.proto

/*
Package register_provider_pb is a generated protocol buffer package.

It is generated from these files:
	provider_register.proto

It has these top-level messages:
	GetPublicKeyReq
	GetPublicKeyResp
	RegisterReq
	RegisterResp
	VerifyBillEmailReq
	VerifyBillEmailResp
	ResendVerifyCodeReq
	ResendVerifyCodeResp
	GetTrackerServerReq
	GetTrackerServerResp
	TrackerServer
*/
package register_provider_pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetPublicKeyReq struct {
	Version uint32 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
}

func (m *GetPublicKeyReq) Reset()                    { *m = GetPublicKeyReq{} }
func (m *GetPublicKeyReq) String() string            { return proto.CompactTextString(m) }
func (*GetPublicKeyReq) ProtoMessage()               {}
func (*GetPublicKeyReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetPublicKeyReq) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

type GetPublicKeyResp struct {
	PublicKey []byte `protobuf:"bytes,1,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	Ip        string `protobuf:"bytes,2,opt,name=ip" json:"ip,omitempty"`
}

func (m *GetPublicKeyResp) Reset()                    { *m = GetPublicKeyResp{} }
func (m *GetPublicKeyResp) String() string            { return proto.CompactTextString(m) }
func (*GetPublicKeyResp) ProtoMessage()               {}
func (*GetPublicKeyResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetPublicKeyResp) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *GetPublicKeyResp) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type RegisterReq struct {
	Version            uint32   `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Timestamp          uint64   `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
	NodeIdEnc          []byte   `protobuf:"bytes,3,opt,name=nodeIdEnc,proto3" json:"nodeIdEnc,omitempty"`
	PublicKeyEnc       []byte   `protobuf:"bytes,4,opt,name=publicKeyEnc,proto3" json:"publicKeyEnc,omitempty"`
	EncryptKeyEnc      []byte   `protobuf:"bytes,5,opt,name=encryptKeyEnc,proto3" json:"encryptKeyEnc,omitempty"`
	WalletAddressEnc   []byte   `protobuf:"bytes,6,opt,name=walletAddressEnc,proto3" json:"walletAddressEnc,omitempty"`
	BillEmailEnc       []byte   `protobuf:"bytes,7,opt,name=billEmailEnc,proto3" json:"billEmailEnc,omitempty"`
	MainStorageVolume  uint64   `protobuf:"varint,8,opt,name=mainStorageVolume" json:"mainStorageVolume,omitempty"`
	UpBandwidth        uint64   `protobuf:"varint,9,opt,name=upBandwidth" json:"upBandwidth,omitempty"`
	DownBandwidth      uint64   `protobuf:"varint,10,opt,name=downBandwidth" json:"downBandwidth,omitempty"`
	TestUpBandwidth    uint64   `protobuf:"varint,11,opt,name=testUpBandwidth" json:"testUpBandwidth,omitempty"`
	TestDownBandwidth  uint64   `protobuf:"varint,12,opt,name=testDownBandwidth" json:"testDownBandwidth,omitempty"`
	Availability       float64  `protobuf:"fixed64,13,opt,name=availability" json:"availability,omitempty"`
	Port               uint32   `protobuf:"varint,14,opt,name=port" json:"port,omitempty"`
	HostEnc            []byte   `protobuf:"bytes,15,opt,name=hostEnc,proto3" json:"hostEnc,omitempty"`
	DynamicDomainEnc   []byte   `protobuf:"bytes,16,opt,name=dynamicDomainEnc,proto3" json:"dynamicDomainEnc,omitempty"`
	ExtraStorageVolume []uint64 `protobuf:"varint,17,rep,packed,name=extraStorageVolume" json:"extraStorageVolume,omitempty"`
	Sign               []byte   `protobuf:"bytes,18,opt,name=sign,proto3" json:"sign,omitempty"`
}

func (m *RegisterReq) Reset()                    { *m = RegisterReq{} }
func (m *RegisterReq) String() string            { return proto.CompactTextString(m) }
func (*RegisterReq) ProtoMessage()               {}
func (*RegisterReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RegisterReq) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *RegisterReq) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *RegisterReq) GetNodeIdEnc() []byte {
	if m != nil {
		return m.NodeIdEnc
	}
	return nil
}

func (m *RegisterReq) GetPublicKeyEnc() []byte {
	if m != nil {
		return m.PublicKeyEnc
	}
	return nil
}

func (m *RegisterReq) GetEncryptKeyEnc() []byte {
	if m != nil {
		return m.EncryptKeyEnc
	}
	return nil
}

func (m *RegisterReq) GetWalletAddressEnc() []byte {
	if m != nil {
		return m.WalletAddressEnc
	}
	return nil
}

func (m *RegisterReq) GetBillEmailEnc() []byte {
	if m != nil {
		return m.BillEmailEnc
	}
	return nil
}

func (m *RegisterReq) GetMainStorageVolume() uint64 {
	if m != nil {
		return m.MainStorageVolume
	}
	return 0
}

func (m *RegisterReq) GetUpBandwidth() uint64 {
	if m != nil {
		return m.UpBandwidth
	}
	return 0
}

func (m *RegisterReq) GetDownBandwidth() uint64 {
	if m != nil {
		return m.DownBandwidth
	}
	return 0
}

func (m *RegisterReq) GetTestUpBandwidth() uint64 {
	if m != nil {
		return m.TestUpBandwidth
	}
	return 0
}

func (m *RegisterReq) GetTestDownBandwidth() uint64 {
	if m != nil {
		return m.TestDownBandwidth
	}
	return 0
}

func (m *RegisterReq) GetAvailability() float64 {
	if m != nil {
		return m.Availability
	}
	return 0
}

func (m *RegisterReq) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *RegisterReq) GetHostEnc() []byte {
	if m != nil {
		return m.HostEnc
	}
	return nil
}

func (m *RegisterReq) GetDynamicDomainEnc() []byte {
	if m != nil {
		return m.DynamicDomainEnc
	}
	return nil
}

func (m *RegisterReq) GetExtraStorageVolume() []uint64 {
	if m != nil {
		return m.ExtraStorageVolume
	}
	return nil
}

func (m *RegisterReq) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

type RegisterResp struct {
	Code   uint32 `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	ErrMsg string `protobuf:"bytes,2,opt,name=errMsg" json:"errMsg,omitempty"`
}

func (m *RegisterResp) Reset()                    { *m = RegisterResp{} }
func (m *RegisterResp) String() string            { return proto.CompactTextString(m) }
func (*RegisterResp) ProtoMessage()               {}
func (*RegisterResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RegisterResp) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *RegisterResp) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type VerifyBillEmailReq struct {
	Version    uint32 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	NodeId     []byte `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Timestamp  uint64 `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	VerifyCode string `protobuf:"bytes,4,opt,name=verifyCode" json:"verifyCode,omitempty"`
	Sign       []byte `protobuf:"bytes,5,opt,name=sign,proto3" json:"sign,omitempty"`
}

func (m *VerifyBillEmailReq) Reset()                    { *m = VerifyBillEmailReq{} }
func (m *VerifyBillEmailReq) String() string            { return proto.CompactTextString(m) }
func (*VerifyBillEmailReq) ProtoMessage()               {}
func (*VerifyBillEmailReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *VerifyBillEmailReq) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *VerifyBillEmailReq) GetNodeId() []byte {
	if m != nil {
		return m.NodeId
	}
	return nil
}

func (m *VerifyBillEmailReq) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *VerifyBillEmailReq) GetVerifyCode() string {
	if m != nil {
		return m.VerifyCode
	}
	return ""
}

func (m *VerifyBillEmailReq) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

type VerifyBillEmailResp struct {
	Code   uint32 `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	ErrMsg string `protobuf:"bytes,2,opt,name=errMsg" json:"errMsg,omitempty"`
}

func (m *VerifyBillEmailResp) Reset()                    { *m = VerifyBillEmailResp{} }
func (m *VerifyBillEmailResp) String() string            { return proto.CompactTextString(m) }
func (*VerifyBillEmailResp) ProtoMessage()               {}
func (*VerifyBillEmailResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *VerifyBillEmailResp) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *VerifyBillEmailResp) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type ResendVerifyCodeReq struct {
	Version   uint32 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	NodeId    []byte `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Timestamp uint64 `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	Sign      []byte `protobuf:"bytes,4,opt,name=sign,proto3" json:"sign,omitempty"`
}

func (m *ResendVerifyCodeReq) Reset()                    { *m = ResendVerifyCodeReq{} }
func (m *ResendVerifyCodeReq) String() string            { return proto.CompactTextString(m) }
func (*ResendVerifyCodeReq) ProtoMessage()               {}
func (*ResendVerifyCodeReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ResendVerifyCodeReq) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ResendVerifyCodeReq) GetNodeId() []byte {
	if m != nil {
		return m.NodeId
	}
	return nil
}

func (m *ResendVerifyCodeReq) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *ResendVerifyCodeReq) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

type ResendVerifyCodeResp struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *ResendVerifyCodeResp) Reset()                    { *m = ResendVerifyCodeResp{} }
func (m *ResendVerifyCodeResp) String() string            { return proto.CompactTextString(m) }
func (*ResendVerifyCodeResp) ProtoMessage()               {}
func (*ResendVerifyCodeResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ResendVerifyCodeResp) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type GetTrackerServerReq struct {
	Version   uint32 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	NodeId    []byte `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Timestamp uint64 `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	Sign      []byte `protobuf:"bytes,4,opt,name=sign,proto3" json:"sign,omitempty"`
}

func (m *GetTrackerServerReq) Reset()                    { *m = GetTrackerServerReq{} }
func (m *GetTrackerServerReq) String() string            { return proto.CompactTextString(m) }
func (*GetTrackerServerReq) ProtoMessage()               {}
func (*GetTrackerServerReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *GetTrackerServerReq) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *GetTrackerServerReq) GetNodeId() []byte {
	if m != nil {
		return m.NodeId
	}
	return nil
}

func (m *GetTrackerServerReq) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *GetTrackerServerReq) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

type GetTrackerServerResp struct {
	Server []*TrackerServer `protobuf:"bytes,1,rep,name=server" json:"server,omitempty"`
}

func (m *GetTrackerServerResp) Reset()                    { *m = GetTrackerServerResp{} }
func (m *GetTrackerServerResp) String() string            { return proto.CompactTextString(m) }
func (*GetTrackerServerResp) ProtoMessage()               {}
func (*GetTrackerServerResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *GetTrackerServerResp) GetServer() []*TrackerServer {
	if m != nil {
		return m.Server
	}
	return nil
}

type TrackerServer struct {
	Server string `protobuf:"bytes,1,opt,name=server" json:"server,omitempty"`
	Port   uint32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
}

func (m *TrackerServer) Reset()                    { *m = TrackerServer{} }
func (m *TrackerServer) String() string            { return proto.CompactTextString(m) }
func (*TrackerServer) ProtoMessage()               {}
func (*TrackerServer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *TrackerServer) GetServer() string {
	if m != nil {
		return m.Server
	}
	return ""
}

func (m *TrackerServer) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func init() {
	proto.RegisterType((*GetPublicKeyReq)(nil), "register_provider_pb.GetPublicKeyReq")
	proto.RegisterType((*GetPublicKeyResp)(nil), "register_provider_pb.GetPublicKeyResp")
	proto.RegisterType((*RegisterReq)(nil), "register_provider_pb.RegisterReq")
	proto.RegisterType((*RegisterResp)(nil), "register_provider_pb.RegisterResp")
	proto.RegisterType((*VerifyBillEmailReq)(nil), "register_provider_pb.VerifyBillEmailReq")
	proto.RegisterType((*VerifyBillEmailResp)(nil), "register_provider_pb.VerifyBillEmailResp")
	proto.RegisterType((*ResendVerifyCodeReq)(nil), "register_provider_pb.ResendVerifyCodeReq")
	proto.RegisterType((*ResendVerifyCodeResp)(nil), "register_provider_pb.ResendVerifyCodeResp")
	proto.RegisterType((*GetTrackerServerReq)(nil), "register_provider_pb.GetTrackerServerReq")
	proto.RegisterType((*GetTrackerServerResp)(nil), "register_provider_pb.GetTrackerServerResp")
	proto.RegisterType((*TrackerServer)(nil), "register_provider_pb.TrackerServer")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ProviderRegisterService service

type ProviderRegisterServiceClient interface {
	GetPublicKey(ctx context.Context, in *GetPublicKeyReq, opts ...grpc.CallOption) (*GetPublicKeyResp, error)
	Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error)
	VerifyBillEmail(ctx context.Context, in *VerifyBillEmailReq, opts ...grpc.CallOption) (*VerifyBillEmailResp, error)
	ResendVerifyCode(ctx context.Context, in *ResendVerifyCodeReq, opts ...grpc.CallOption) (*ResendVerifyCodeResp, error)
	GetTrackerServer(ctx context.Context, in *GetTrackerServerReq, opts ...grpc.CallOption) (*GetTrackerServerResp, error)
}

type providerRegisterServiceClient struct {
	cc *grpc.ClientConn
}

func NewProviderRegisterServiceClient(cc *grpc.ClientConn) ProviderRegisterServiceClient {
	return &providerRegisterServiceClient{cc}
}

func (c *providerRegisterServiceClient) GetPublicKey(ctx context.Context, in *GetPublicKeyReq, opts ...grpc.CallOption) (*GetPublicKeyResp, error) {
	out := new(GetPublicKeyResp)
	err := grpc.Invoke(ctx, "/register_provider_pb.ProviderRegisterService/GetPublicKey", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerRegisterServiceClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error) {
	out := new(RegisterResp)
	err := grpc.Invoke(ctx, "/register_provider_pb.ProviderRegisterService/Register", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerRegisterServiceClient) VerifyBillEmail(ctx context.Context, in *VerifyBillEmailReq, opts ...grpc.CallOption) (*VerifyBillEmailResp, error) {
	out := new(VerifyBillEmailResp)
	err := grpc.Invoke(ctx, "/register_provider_pb.ProviderRegisterService/VerifyBillEmail", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerRegisterServiceClient) ResendVerifyCode(ctx context.Context, in *ResendVerifyCodeReq, opts ...grpc.CallOption) (*ResendVerifyCodeResp, error) {
	out := new(ResendVerifyCodeResp)
	err := grpc.Invoke(ctx, "/register_provider_pb.ProviderRegisterService/ResendVerifyCode", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerRegisterServiceClient) GetTrackerServer(ctx context.Context, in *GetTrackerServerReq, opts ...grpc.CallOption) (*GetTrackerServerResp, error) {
	out := new(GetTrackerServerResp)
	err := grpc.Invoke(ctx, "/register_provider_pb.ProviderRegisterService/GetTrackerServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ProviderRegisterService service

type ProviderRegisterServiceServer interface {
	GetPublicKey(context.Context, *GetPublicKeyReq) (*GetPublicKeyResp, error)
	Register(context.Context, *RegisterReq) (*RegisterResp, error)
	VerifyBillEmail(context.Context, *VerifyBillEmailReq) (*VerifyBillEmailResp, error)
	ResendVerifyCode(context.Context, *ResendVerifyCodeReq) (*ResendVerifyCodeResp, error)
	GetTrackerServer(context.Context, *GetTrackerServerReq) (*GetTrackerServerResp, error)
}

func RegisterProviderRegisterServiceServer(s *grpc.Server, srv ProviderRegisterServiceServer) {
	s.RegisterService(&_ProviderRegisterService_serviceDesc, srv)
}

func _ProviderRegisterService_GetPublicKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPublicKeyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegisterServiceServer).GetPublicKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/register_provider_pb.ProviderRegisterService/GetPublicKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegisterServiceServer).GetPublicKey(ctx, req.(*GetPublicKeyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderRegisterService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegisterServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/register_provider_pb.ProviderRegisterService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegisterServiceServer).Register(ctx, req.(*RegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderRegisterService_VerifyBillEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyBillEmailReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegisterServiceServer).VerifyBillEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/register_provider_pb.ProviderRegisterService/VerifyBillEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegisterServiceServer).VerifyBillEmail(ctx, req.(*VerifyBillEmailReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderRegisterService_ResendVerifyCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResendVerifyCodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegisterServiceServer).ResendVerifyCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/register_provider_pb.ProviderRegisterService/ResendVerifyCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegisterServiceServer).ResendVerifyCode(ctx, req.(*ResendVerifyCodeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderRegisterService_GetTrackerServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTrackerServerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderRegisterServiceServer).GetTrackerServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/register_provider_pb.ProviderRegisterService/GetTrackerServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderRegisterServiceServer).GetTrackerServer(ctx, req.(*GetTrackerServerReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProviderRegisterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "register_provider_pb.ProviderRegisterService",
	HandlerType: (*ProviderRegisterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPublicKey",
			Handler:    _ProviderRegisterService_GetPublicKey_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _ProviderRegisterService_Register_Handler,
		},
		{
			MethodName: "VerifyBillEmail",
			Handler:    _ProviderRegisterService_VerifyBillEmail_Handler,
		},
		{
			MethodName: "ResendVerifyCode",
			Handler:    _ProviderRegisterService_ResendVerifyCode_Handler,
		},
		{
			MethodName: "GetTrackerServer",
			Handler:    _ProviderRegisterService_GetTrackerServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "provider_register.proto",
}

func init() { proto.RegisterFile("provider_register.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 702 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0xdb, 0x4f, 0xdb, 0x3e,
	0x18, 0x25, 0x6d, 0x28, 0xf0, 0xb5, 0xa5, 0xc5, 0x20, 0x88, 0xaa, 0x9f, 0x7e, 0xea, 0xb2, 0x8b,
	0x02, 0x9b, 0xaa, 0x89, 0xbd, 0x8d, 0x97, 0xc1, 0x40, 0x68, 0x9a, 0x26, 0xa1, 0x74, 0xe3, 0x15,
	0xa5, 0x89, 0x57, 0x2c, 0xd2, 0xc4, 0xb3, 0xdd, 0xb2, 0xfc, 0x11, 0x7b, 0xda, 0xcb, 0xfe, 0xdc,
	0xc9, 0xce, 0xa5, 0xb9, 0x71, 0x99, 0xb4, 0xbd, 0xe5, 0x3b, 0x3e, 0xfe, 0xce, 0xb1, 0x7d, 0xec,
	0x16, 0xf6, 0x28, 0x0b, 0x17, 0xc4, 0xc3, 0xec, 0x8a, 0xe1, 0x29, 0xe1, 0x02, 0xb3, 0x11, 0x65,
	0xa1, 0x08, 0xd1, 0x4e, 0x5a, 0x5f, 0x65, 0x0c, 0x3a, 0x31, 0x5f, 0x42, 0xef, 0x1c, 0x8b, 0x8b,
	0xf9, 0xc4, 0x27, 0xee, 0x47, 0x1c, 0xd9, 0xf8, 0x1b, 0x32, 0x60, 0x6d, 0x81, 0x19, 0x27, 0x61,
	0x60, 0x68, 0x43, 0xcd, 0xea, 0xda, 0x69, 0x69, 0xbe, 0x83, 0x7e, 0x91, 0xcc, 0x29, 0xfa, 0x0f,
	0x36, 0x68, 0x0a, 0x28, 0x7e, 0xc7, 0x5e, 0x02, 0x68, 0x13, 0x1a, 0x84, 0x1a, 0x8d, 0xa1, 0x66,
	0x6d, 0xd8, 0x0d, 0x42, 0xcd, 0x1f, 0xab, 0xd0, 0xb6, 0x13, 0x1f, 0xf7, 0x6a, 0xc9, 0xbe, 0x82,
	0xcc, 0x30, 0x17, 0xce, 0x2c, 0x6e, 0xa0, 0xdb, 0x4b, 0x40, 0x8e, 0x06, 0xa1, 0x87, 0x3f, 0x78,
	0x67, 0x81, 0x6b, 0x34, 0x63, 0xd5, 0x0c, 0x40, 0x26, 0x74, 0x32, 0x0b, 0x92, 0xa0, 0x2b, 0x42,
	0x01, 0x43, 0xcf, 0xa0, 0x8b, 0x03, 0x97, 0x45, 0x54, 0x24, 0xa4, 0x55, 0x45, 0x2a, 0x82, 0xe8,
	0x00, 0xfa, 0xb7, 0x8e, 0xef, 0x63, 0x71, 0xec, 0x79, 0x0c, 0x73, 0x2e, 0x89, 0x2d, 0x45, 0xac,
	0xe0, 0x52, 0x75, 0x42, 0x7c, 0xff, 0x6c, 0xe6, 0x10, 0x5f, 0xf2, 0xd6, 0x62, 0xd5, 0x3c, 0x86,
	0x5e, 0xc1, 0xd6, 0xcc, 0x21, 0xc1, 0x58, 0x84, 0xcc, 0x99, 0xe2, 0xcb, 0xd0, 0x9f, 0xcf, 0xb0,
	0xb1, 0xae, 0x56, 0x57, 0x1d, 0x40, 0x43, 0x68, 0xcf, 0xe9, 0x89, 0x13, 0x78, 0xb7, 0xc4, 0x13,
	0xd7, 0xc6, 0x86, 0xe2, 0xe5, 0x21, 0xb9, 0x0a, 0x2f, 0xbc, 0x0d, 0x96, 0x1c, 0x50, 0x9c, 0x22,
	0x88, 0x2c, 0xe8, 0x09, 0xcc, 0xc5, 0x97, 0x5c, 0xaf, 0xb6, 0xe2, 0x95, 0x61, 0xe9, 0x4f, 0x42,
	0xa7, 0x85, 0x9e, 0x9d, 0xd8, 0x5f, 0x65, 0x40, 0xae, 0xd8, 0x59, 0x38, 0xc4, 0x77, 0x26, 0xc4,
	0x27, 0x22, 0x32, 0xba, 0x43, 0xcd, 0xd2, 0xec, 0x02, 0x86, 0x10, 0xe8, 0x34, 0x64, 0xc2, 0xd8,
	0x54, 0xc7, 0xab, 0xbe, 0xe5, 0xa9, 0x5f, 0x87, 0x5c, 0xc8, 0x4d, 0xea, 0xa9, 0x4d, 0x4a, 0x4b,
	0xb9, 0xdf, 0x5e, 0x14, 0x38, 0x33, 0xe2, 0x9e, 0x86, 0x72, 0x3f, 0x24, 0xa5, 0x1f, 0xef, 0x77,
	0x19, 0x47, 0x23, 0x40, 0xf8, 0xbb, 0x60, 0x4e, 0x71, 0x33, 0xb7, 0x86, 0x4d, 0x4b, 0xb7, 0x6b,
	0x46, 0xa4, 0x13, 0x4e, 0xa6, 0x81, 0x81, 0x54, 0x3f, 0xf5, 0x6d, 0xbe, 0x85, 0xce, 0x32, 0x8e,
	0x9c, 0x4a, 0x8e, 0x1b, 0x7a, 0x38, 0x09, 0xa3, 0xfa, 0x46, 0xbb, 0xd0, 0xc2, 0x8c, 0x7d, 0xe2,
	0xd3, 0x24, 0xc7, 0x49, 0x65, 0xfe, 0xd2, 0x00, 0x5d, 0x62, 0x46, 0xbe, 0x46, 0x27, 0xe9, 0x11,
	0xdf, 0x1f, 0xe9, 0x5d, 0x68, 0xc5, 0x19, 0x55, 0x8d, 0x3a, 0x76, 0x52, 0x15, 0xa3, 0xde, 0x2c,
	0x47, 0xfd, 0x7f, 0x80, 0x85, 0x52, 0x79, 0x2f, 0x8d, 0xe9, 0xca, 0x42, 0x0e, 0xc9, 0x96, 0xb5,
	0x9a, 0x5b, 0xd6, 0x31, 0x6c, 0x57, 0x9c, 0xfd, 0xe1, 0xea, 0x22, 0xd8, 0xb6, 0x31, 0xc7, 0x81,
	0x77, 0x99, 0x49, 0xfd, 0x8b, 0xd5, 0xa5, 0xee, 0xf5, 0x9c, 0xfb, 0xd7, 0xb0, 0x53, 0x95, 0xe6,
	0x54, 0x6a, 0xf3, 0xb9, 0xeb, 0x62, 0xce, 0x95, 0xf6, 0xba, 0x9d, 0x96, 0xd2, 0xec, 0x39, 0x16,
	0x9f, 0x99, 0xe3, 0xde, 0x60, 0x36, 0xc6, 0x6c, 0xf1, 0xd0, 0xeb, 0xf2, 0xf7, 0xcc, 0x8e, 0x61,
	0xa7, 0x2a, 0xcd, 0x29, 0x3a, 0x82, 0x16, 0x57, 0x95, 0xa1, 0x0d, 0x9b, 0x56, 0xfb, 0xf0, 0xe9,
	0xa8, 0xee, 0xfd, 0x1d, 0x15, 0x27, 0x26, 0x53, 0xcc, 0x23, 0xe8, 0x16, 0x06, 0xa4, 0xdf, 0xac,
	0x9b, 0x3a, 0xa5, 0xb8, 0xca, 0x6e, 0x57, 0x63, 0x79, 0xbb, 0x0e, 0x7f, 0xea, 0xb0, 0x77, 0x91,
	0x48, 0xa4, 0xe1, 0x96, 0x6d, 0x88, 0x8b, 0xd1, 0x15, 0x74, 0xf2, 0x2f, 0x38, 0x7a, 0x5e, 0xef,
	0xaa, 0xf4, 0x93, 0x30, 0x78, 0xf1, 0x18, 0x1a, 0xa7, 0xe6, 0x0a, 0x1a, 0xc3, 0x7a, 0xaa, 0x89,
	0x9e, 0xd4, 0xcf, 0xca, 0xbd, 0xff, 0x03, 0xf3, 0x21, 0x8a, 0x6a, 0x7a, 0x0d, 0xbd, 0x52, 0x9c,
	0x91, 0x55, 0x3f, 0xb1, 0x7a, 0x1f, 0x07, 0xfb, 0x8f, 0x64, 0x2a, 0xa5, 0x1b, 0xe8, 0x97, 0xa3,
	0x87, 0xf6, 0xef, 0xf2, 0x58, 0xb9, 0x1d, 0x83, 0x83, 0xc7, 0x52, 0x53, 0xb1, 0x72, 0x74, 0xee,
	0x12, 0xab, 0x49, 0xf7, 0x5d, 0x62, 0x75, 0x69, 0x34, 0x57, 0x26, 0x2d, 0xf5, 0x2f, 0xe0, 0xcd,
	0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdd, 0x7e, 0x99, 0x63, 0x20, 0x08, 0x00, 0x00,
}
