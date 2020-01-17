// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package jobber

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

type Response_ResponseStatus int32

const (
	Response_StillProcessing      Response_ResponseStatus = 0
	Response_FinishedSuccessfully Response_ResponseStatus = 1
	Response_FinishedWithError    Response_ResponseStatus = 2
)

var Response_ResponseStatus_name = map[int32]string{
	0: "StillProcessing",
	1: "FinishedSuccessfully",
	2: "FinishedWithError",
}

var Response_ResponseStatus_value = map[string]int32{
	"StillProcessing":      0,
	"FinishedSuccessfully": 1,
	"FinishedWithError":    2,
}

func (x Response_ResponseStatus) String() string {
	return proto.EnumName(Response_ResponseStatus_name, int32(x))
}

func (Response_ResponseStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2, 0}
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Instruction struct {
	Instruction          []byte   `protobuf:"bytes,1,opt,name=Instruction,proto3" json:"Instruction,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Instruction) Reset()         { *m = Instruction{} }
func (m *Instruction) String() string { return proto.CompactTextString(m) }
func (*Instruction) ProtoMessage()    {}
func (*Instruction) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *Instruction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Instruction.Unmarshal(m, b)
}
func (m *Instruction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Instruction.Marshal(b, m, deterministic)
}
func (m *Instruction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Instruction.Merge(m, src)
}
func (m *Instruction) XXX_Size() int {
	return xxx_messageInfo_Instruction.Size(m)
}
func (m *Instruction) XXX_DiscardUnknown() {
	xxx_messageInfo_Instruction.DiscardUnknown(m)
}

var xxx_messageInfo_Instruction proto.InternalMessageInfo

func (m *Instruction) GetInstruction() []byte {
	if m != nil {
		return m.Instruction
	}
	return nil
}

type Response struct {
	Status               Response_ResponseStatus `protobuf:"varint,1,opt,name=Status,proto3,enum=jobber.Response_ResponseStatus" json:"Status,omitempty"`
	Message              []byte                  `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
	Error                string                  `protobuf:"bytes,3,opt,name=Error,proto3" json:"Error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetStatus() Response_ResponseStatus {
	if m != nil {
		return m.Status
	}
	return Response_StillProcessing
}

func (m *Response) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterEnum("jobber.Response_ResponseStatus", Response_ResponseStatus_name, Response_ResponseStatus_value)
	proto.RegisterType((*Empty)(nil), "jobber.Empty")
	proto.RegisterType((*Instruction)(nil), "jobber.Instruction")
	proto.RegisterType((*Response)(nil), "jobber.Response")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 272 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0xcd, 0xa4, 0xad, 0xbb, 0xba, 0x59, 0xef, 0x26, 0x94, 0xbd, 0x58, 0x0a, 0xc2, 0x9e,
	0xaa, 0xd4, 0x07, 0x7f, 0xc1, 0x44, 0x05, 0x61, 0xb4, 0xa0, 0xcf, 0x6d, 0xbd, 0xdb, 0x32, 0xb2,
	0xa6, 0x24, 0xa9, 0xb0, 0xdf, 0xe8, 0x9f, 0x12, 0x53, 0x2b, 0xeb, 0xde, 0x72, 0xee, 0xf9, 0x38,
	0x49, 0xce, 0x85, 0x61, 0x5e, 0xf3, 0xb8, 0x56, 0xd2, 0x48, 0x74, 0xb7, 0xb2, 0x28, 0x48, 0x45,
	0x1e, 0x38, 0x8b, 0x5d, 0x6d, 0xf6, 0xd1, 0x1d, 0x9c, 0xbf, 0x54, 0xda, 0xa8, 0xa6, 0x34, 0x5c,
	0x56, 0x18, 0xf6, 0x64, 0xc0, 0x42, 0x36, 0xbf, 0x48, 0x0f, 0x47, 0xd1, 0x37, 0x83, 0xb3, 0x94,
	0x74, 0x2d, 0x2b, 0x4d, 0xf8, 0x08, 0x6e, 0x66, 0x72, 0xd3, 0x68, 0x4b, 0x8e, 0x93, 0x9b, 0xb8,
	0xcd, 0x8f, 0x3b, 0xe2, 0xff, 0xd0, 0x62, 0xe9, 0x1f, 0x8e, 0x01, 0x78, 0x6f, 0xa4, 0x75, 0xbe,
	0xa6, 0x60, 0x60, 0xef, 0xe8, 0x24, 0x4e, 0xc1, 0x59, 0x28, 0x25, 0x55, 0x70, 0x1a, 0xb2, 0xf9,
	0x30, 0x6d, 0x45, 0xf4, 0x0e, 0xe3, 0x7e, 0x12, 0x4e, 0xe0, 0x32, 0x33, 0x5c, 0x88, 0xa5, 0x92,
	0x25, 0x69, 0xcd, 0xab, 0xb5, 0x7f, 0x82, 0x01, 0x4c, 0x9f, 0x78, 0xc5, 0xf5, 0x86, 0x3e, 0xb3,
	0xa6, 0xfc, 0x9d, 0xaf, 0x1a, 0x21, 0xf6, 0x3e, 0xc3, 0x6b, 0xb8, 0xea, 0x9c, 0x0f, 0x6e, 0x36,
	0x36, 0xd5, 0x1f, 0x24, 0x5b, 0x18, 0xbd, 0xda, 0x17, 0x67, 0xa4, 0xbe, 0x78, 0x49, 0x78, 0x0b,
	0xce, 0x33, 0x09, 0x21, 0x71, 0xd4, 0x7d, 0xc5, 0xf6, 0x34, 0xeb, 0x4b, 0x4c, 0xc0, 0x5b, 0x92,
	0x5a, 0x49, 0xb5, 0xc3, 0x49, 0xe7, 0x1c, 0xb4, 0x34, 0xf3, 0x8f, 0x8b, 0xb8, 0x67, 0x85, 0x6b,
	0x57, 0xf0, 0xf0, 0x13, 0x00, 0x00, 0xff, 0xff, 0x8b, 0x1c, 0x02, 0x5e, 0x8f, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JobberServiceClient is the client API for JobberService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JobberServiceClient interface {
	Hello(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Perform(ctx context.Context, in *Instruction, opts ...grpc.CallOption) (JobberService_PerformClient, error)
}

type jobberServiceClient struct {
	cc *grpc.ClientConn
}

func NewJobberServiceClient(cc *grpc.ClientConn) JobberServiceClient {
	return &jobberServiceClient{cc}
}

func (c *jobberServiceClient) Hello(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/jobber.JobberService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobberServiceClient) Perform(ctx context.Context, in *Instruction, opts ...grpc.CallOption) (JobberService_PerformClient, error) {
	stream, err := c.cc.NewStream(ctx, &_JobberService_serviceDesc.Streams[0], "/jobber.JobberService/Perform", opts...)
	if err != nil {
		return nil, err
	}
	x := &jobberServicePerformClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type JobberService_PerformClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type jobberServicePerformClient struct {
	grpc.ClientStream
}

func (x *jobberServicePerformClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JobberServiceServer is the server API for JobberService service.
type JobberServiceServer interface {
	Hello(context.Context, *Empty) (*Empty, error)
	Perform(*Instruction, JobberService_PerformServer) error
}

// UnimplementedJobberServiceServer can be embedded to have forward compatible implementations.
type UnimplementedJobberServiceServer struct {
}

func (*UnimplementedJobberServiceServer) Hello(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (*UnimplementedJobberServiceServer) Perform(req *Instruction, srv JobberService_PerformServer) error {
	return status.Errorf(codes.Unimplemented, "method Perform not implemented")
}

func RegisterJobberServiceServer(s *grpc.Server, srv JobberServiceServer) {
	s.RegisterService(&_JobberService_serviceDesc, srv)
}

func _JobberService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobberServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobber.JobberService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobberServiceServer).Hello(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobberService_Perform_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Instruction)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JobberServiceServer).Perform(m, &jobberServicePerformServer{stream})
}

type JobberService_PerformServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type jobberServicePerformServer struct {
	grpc.ServerStream
}

func (x *jobberServicePerformServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

var _JobberService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "jobber.JobberService",
	HandlerType: (*JobberServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _JobberService_Hello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Perform",
			Handler:       _JobberService_Perform_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api.proto",
}