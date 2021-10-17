// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package pb // import "./"

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

type GetPostReq struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPostReq) Reset()         { *m = GetPostReq{} }
func (m *GetPostReq) String() string { return proto.CompactTextString(m) }
func (*GetPostReq) ProtoMessage()    {}
func (*GetPostReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_post_53c430a5a9bb96a8, []int{0}
}
func (m *GetPostReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPostReq.Unmarshal(m, b)
}
func (m *GetPostReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPostReq.Marshal(b, m, deterministic)
}
func (dst *GetPostReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPostReq.Merge(dst, src)
}
func (m *GetPostReq) XXX_Size() int {
	return xxx_messageInfo_GetPostReq.Size(m)
}
func (m *GetPostReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPostReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetPostReq proto.InternalMessageInfo

func (m *GetPostReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// HelloReply 响应数据格式
type GetPostResp struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Text                 string   `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPostResp) Reset()         { *m = GetPostResp{} }
func (m *GetPostResp) String() string { return proto.CompactTextString(m) }
func (*GetPostResp) ProtoMessage()    {}
func (*GetPostResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_post_53c430a5a9bb96a8, []int{1}
}
func (m *GetPostResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPostResp.Unmarshal(m, b)
}
func (m *GetPostResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPostResp.Marshal(b, m, deterministic)
}
func (dst *GetPostResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPostResp.Merge(dst, src)
}
func (m *GetPostResp) XXX_Size() int {
	return xxx_messageInfo_GetPostResp.Size(m)
}
func (m *GetPostResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPostResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetPostResp proto.InternalMessageInfo

func (m *GetPostResp) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *GetPostResp) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func init() {
	proto.RegisterType((*GetPostReq)(nil), "GetPostReq")
	proto.RegisterType((*GetPostResp)(nil), "GetPostResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PostServiceClient is the export API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PostServiceClient interface {
	GetPost(ctx context.Context, in *GetPostReq, opts ...grpc.CallOption) (*GetPostResp, error)
}

type postServiceClient struct {
	cc *grpc.ClientConn
}

func NewPostServiceClient(cc *grpc.ClientConn) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) GetPost(ctx context.Context, in *GetPostReq, opts ...grpc.CallOption) (*GetPostResp, error) {
	out := new(GetPostResp)
	err := c.cc.Invoke(ctx, "/PostService/GetPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
type PostServiceServer interface {
	GetPost(context.Context, *GetPostReq) (*GetPostResp, error)
}

func RegisterPostServiceServer(s *grpc.Server, srv PostServiceServer) {
	s.RegisterService(&_PostService_serviceDesc, srv)
}

func _PostService_GetPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PostService/GetPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetPost(ctx, req.(*GetPostReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _PostService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPost",
			Handler:    _PostService_GetPost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}

func init() { proto.RegisterFile("user.proto", fileDescriptor_post_53c430a5a9bb96a8) }

var fileDescriptor_post_53c430a5a9bb96a8 = []byte{
	// 145 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0xc8, 0x2f, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x92, 0xe1, 0xe2, 0x72, 0x4f, 0x2d, 0x09, 0xc8, 0x2f,
	0x2e, 0x09, 0x4a, 0x2d, 0x14, 0xe2, 0xe3, 0x62, 0xf2, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0,
	0x0c, 0x62, 0xf2, 0x4c, 0x51, 0x32, 0xe6, 0xe2, 0x86, 0xcb, 0x16, 0x17, 0x08, 0x09, 0x70, 0x31,
	0x97, 0x66, 0x42, 0xe4, 0x99, 0x83, 0x40, 0x4c, 0x21, 0x21, 0x2e, 0x96, 0x92, 0xd4, 0x8a, 0x12,
	0x09, 0x26, 0xb0, 0x16, 0x30, 0xdb, 0xc8, 0x94, 0x8b, 0x1b, 0xa4, 0x23, 0x38, 0xb5, 0xa8, 0x2c,
	0x33, 0x39, 0x55, 0x48, 0x8d, 0x8b, 0x1d, 0x6a, 0x86, 0x10, 0xb7, 0x1e, 0xc2, 0x2e, 0x29, 0x1e,
	0x3d, 0x24, 0xa3, 0x95, 0x18, 0x9c, 0xd8, 0xa3, 0x58, 0xf5, 0xf4, 0xad, 0x0b, 0x92, 0x92, 0xd8,
	0xc0, 0x2e, 0x33, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x02, 0x4d, 0x52, 0xac, 0xa7, 0x00, 0x00,
	0x00,
}