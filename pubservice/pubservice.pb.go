// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pubservice.proto

package pubservice

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type String struct {
	Topic                string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cc011b7714728a2, []int{0}
}

func (m *String) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_String.Unmarshal(m, b)
}
func (m *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_String.Marshal(b, m, deterministic)
}
func (m *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(m, src)
}
func (m *String) XXX_Size() int {
	return xxx_messageInfo_String.Size(m)
}
func (m *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(m)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (m *String) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *String) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*String)(nil), "pubservice.String")
}

func init() {
	proto.RegisterFile("pubservice.proto", fileDescriptor_7cc011b7714728a2)
}

var fileDescriptor_7cc011b7714728a2 = []byte{
	// 156 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x28, 0x4d, 0x2a,
	0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88,
	0x28, 0x99, 0x70, 0xb1, 0x05, 0x97, 0x14, 0x65, 0xe6, 0xa5, 0x0b, 0x89, 0x70, 0xb1, 0x96, 0xe4,
	0x17, 0x64, 0x26, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x38, 0x20, 0xd1, 0xb2, 0xc4,
	0x9c, 0xd2, 0x54, 0x09, 0x26, 0x88, 0x28, 0x98, 0x63, 0xb4, 0x95, 0x91, 0x8b, 0x37, 0xa0, 0x34,
	0xa9, 0xb8, 0x34, 0x29, 0x18, 0x62, 0x8e, 0x90, 0x21, 0x17, 0x7b, 0x40, 0x69, 0x52, 0x4e, 0x66,
	0x71, 0x86, 0x90, 0x90, 0x1e, 0x92, 0x8d, 0x10, 0xc3, 0xa5, 0xb0, 0x88, 0x09, 0x59, 0x71, 0xf1,
	0x05, 0x97, 0x26, 0x15, 0x27, 0x17, 0x65, 0x26, 0xa5, 0x86, 0x80, 0x2d, 0x23, 0x52, 0xa7, 0x01,
	0xa3, 0x90, 0x29, 0x17, 0x27, 0x5c, 0x2f, 0xf1, 0xda, 0x92, 0xd8, 0xc0, 0x01, 0x60, 0x0c, 0x08,
	0x00, 0x00, 0xff, 0xff, 0xef, 0x49, 0x02, 0xef, 0x14, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PubsubServiceClient is the client API for PubsubService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PubsubServiceClient interface {
	Publish(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	SubscribeTopic(ctx context.Context, in *String, opts ...grpc.CallOption) (PubsubService_SubscribeTopicClient, error)
	Subscribe(ctx context.Context, in *String, opts ...grpc.CallOption) (PubsubService_SubscribeClient, error)
}

type pubsubServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPubsubServiceClient(cc grpc.ClientConnInterface) PubsubServiceClient {
	return &pubsubServiceClient{cc}
}

func (c *pubsubServiceClient) Publish(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/pubservice.PubsubService/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubsubServiceClient) SubscribeTopic(ctx context.Context, in *String, opts ...grpc.CallOption) (PubsubService_SubscribeTopicClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PubsubService_serviceDesc.Streams[0], "/pubservice.PubsubService/SubscribeTopic", opts...)
	if err != nil {
		return nil, err
	}
	x := &pubsubServiceSubscribeTopicClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PubsubService_SubscribeTopicClient interface {
	Recv() (*String, error)
	grpc.ClientStream
}

type pubsubServiceSubscribeTopicClient struct {
	grpc.ClientStream
}

func (x *pubsubServiceSubscribeTopicClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pubsubServiceClient) Subscribe(ctx context.Context, in *String, opts ...grpc.CallOption) (PubsubService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PubsubService_serviceDesc.Streams[1], "/pubservice.PubsubService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &pubsubServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PubsubService_SubscribeClient interface {
	Recv() (*String, error)
	grpc.ClientStream
}

type pubsubServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *pubsubServiceSubscribeClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PubsubServiceServer is the server API for PubsubService service.
type PubsubServiceServer interface {
	Publish(context.Context, *String) (*String, error)
	SubscribeTopic(*String, PubsubService_SubscribeTopicServer) error
	Subscribe(*String, PubsubService_SubscribeServer) error
}

// UnimplementedPubsubServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPubsubServiceServer struct {
}

func (*UnimplementedPubsubServiceServer) Publish(ctx context.Context, req *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (*UnimplementedPubsubServiceServer) SubscribeTopic(req *String, srv PubsubService_SubscribeTopicServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeTopic not implemented")
}
func (*UnimplementedPubsubServiceServer) Subscribe(req *String, srv PubsubService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}

func RegisterPubsubServiceServer(s *grpc.Server, srv PubsubServiceServer) {
	s.RegisterService(&_PubsubService_serviceDesc, srv)
}

func _PubsubService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubsubServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pubservice.PubsubService/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubsubServiceServer).Publish(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubsubService_SubscribeTopic_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(String)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PubsubServiceServer).SubscribeTopic(m, &pubsubServiceSubscribeTopicServer{stream})
}

type PubsubService_SubscribeTopicServer interface {
	Send(*String) error
	grpc.ServerStream
}

type pubsubServiceSubscribeTopicServer struct {
	grpc.ServerStream
}

func (x *pubsubServiceSubscribeTopicServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func _PubsubService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(String)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PubsubServiceServer).Subscribe(m, &pubsubServiceSubscribeServer{stream})
}

type PubsubService_SubscribeServer interface {
	Send(*String) error
	grpc.ServerStream
}

type pubsubServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *pubsubServiceSubscribeServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

var _PubsubService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pubservice.PubsubService",
	HandlerType: (*PubsubServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _PubsubService_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeTopic",
			Handler:       _PubsubService_SubscribeTopic_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _PubsubService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pubservice.proto",
}
