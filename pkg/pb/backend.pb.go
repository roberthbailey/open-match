// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/backend.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
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

type FunctionConfig_Type int32

const (
	FunctionConfig_GRPC FunctionConfig_Type = 0
	FunctionConfig_REST FunctionConfig_Type = 1
)

var FunctionConfig_Type_name = map[int32]string{
	0: "GRPC",
	1: "REST",
}

var FunctionConfig_Type_value = map[string]int32{
	"GRPC": 0,
	"REST": 1,
}

func (x FunctionConfig_Type) String() string {
	return proto.EnumName(FunctionConfig_Type_name, int32(x))
}

func (FunctionConfig_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8dab762378f455cd, []int{0, 0}
}

// Configuration for the Match Function to be triggered by Open Match to
// generate proposals.
type FunctionConfig struct {
	Host                 string              `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port                 int32               `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Type                 FunctionConfig_Type `protobuf:"varint,3,opt,name=type,proto3,enum=api.FunctionConfig_Type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *FunctionConfig) Reset()         { *m = FunctionConfig{} }
func (m *FunctionConfig) String() string { return proto.CompactTextString(m) }
func (*FunctionConfig) ProtoMessage()    {}
func (*FunctionConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_8dab762378f455cd, []int{0}
}

func (m *FunctionConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FunctionConfig.Unmarshal(m, b)
}
func (m *FunctionConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FunctionConfig.Marshal(b, m, deterministic)
}
func (m *FunctionConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FunctionConfig.Merge(m, src)
}
func (m *FunctionConfig) XXX_Size() int {
	return xxx_messageInfo_FunctionConfig.Size(m)
}
func (m *FunctionConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_FunctionConfig.DiscardUnknown(m)
}

var xxx_messageInfo_FunctionConfig proto.InternalMessageInfo

func (m *FunctionConfig) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *FunctionConfig) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *FunctionConfig) GetType() FunctionConfig_Type {
	if m != nil {
		return m.Type
	}
	return FunctionConfig_GRPC
}

type FetchMatchesRequest struct {
	// Configuration of the MatchFunction to be executed for the given list of MatchProfiles
	Config *FunctionConfig `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	// MatchProfiles for which this MatchFunction should be executed.
	Profile              []*MatchProfile `protobuf:"bytes,2,rep,name=profile,proto3" json:"profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *FetchMatchesRequest) Reset()         { *m = FetchMatchesRequest{} }
func (m *FetchMatchesRequest) String() string { return proto.CompactTextString(m) }
func (*FetchMatchesRequest) ProtoMessage()    {}
func (*FetchMatchesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8dab762378f455cd, []int{1}
}

func (m *FetchMatchesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FetchMatchesRequest.Unmarshal(m, b)
}
func (m *FetchMatchesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FetchMatchesRequest.Marshal(b, m, deterministic)
}
func (m *FetchMatchesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FetchMatchesRequest.Merge(m, src)
}
func (m *FetchMatchesRequest) XXX_Size() int {
	return xxx_messageInfo_FetchMatchesRequest.Size(m)
}
func (m *FetchMatchesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FetchMatchesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FetchMatchesRequest proto.InternalMessageInfo

func (m *FetchMatchesRequest) GetConfig() *FunctionConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *FetchMatchesRequest) GetProfile() []*MatchProfile {
	if m != nil {
		return m.Profile
	}
	return nil
}

type FetchMatchesResponse struct {
	// Result Match for the requested MatchProfile.
	Match                []*Match `protobuf:"bytes,1,rep,name=match,proto3" json:"match,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FetchMatchesResponse) Reset()         { *m = FetchMatchesResponse{} }
func (m *FetchMatchesResponse) String() string { return proto.CompactTextString(m) }
func (*FetchMatchesResponse) ProtoMessage()    {}
func (*FetchMatchesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8dab762378f455cd, []int{2}
}

func (m *FetchMatchesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FetchMatchesResponse.Unmarshal(m, b)
}
func (m *FetchMatchesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FetchMatchesResponse.Marshal(b, m, deterministic)
}
func (m *FetchMatchesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FetchMatchesResponse.Merge(m, src)
}
func (m *FetchMatchesResponse) XXX_Size() int {
	return xxx_messageInfo_FetchMatchesResponse.Size(m)
}
func (m *FetchMatchesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FetchMatchesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FetchMatchesResponse proto.InternalMessageInfo

func (m *FetchMatchesResponse) GetMatch() []*Match {
	if m != nil {
		return m.Match
	}
	return nil
}

type AssignTicketsRequest struct {
	// List of Ticket IDs for which the Assignment is to be made.
	TicketId []string `protobuf:"bytes,1,rep,name=ticket_id,json=ticketId,proto3" json:"ticket_id,omitempty"`
	// Assignment to be associated with the Ticket IDs.
	Assignment           *Assignment `protobuf:"bytes,2,opt,name=assignment,proto3" json:"assignment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *AssignTicketsRequest) Reset()         { *m = AssignTicketsRequest{} }
func (m *AssignTicketsRequest) String() string { return proto.CompactTextString(m) }
func (*AssignTicketsRequest) ProtoMessage()    {}
func (*AssignTicketsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8dab762378f455cd, []int{3}
}

func (m *AssignTicketsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AssignTicketsRequest.Unmarshal(m, b)
}
func (m *AssignTicketsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AssignTicketsRequest.Marshal(b, m, deterministic)
}
func (m *AssignTicketsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssignTicketsRequest.Merge(m, src)
}
func (m *AssignTicketsRequest) XXX_Size() int {
	return xxx_messageInfo_AssignTicketsRequest.Size(m)
}
func (m *AssignTicketsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AssignTicketsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AssignTicketsRequest proto.InternalMessageInfo

func (m *AssignTicketsRequest) GetTicketId() []string {
	if m != nil {
		return m.TicketId
	}
	return nil
}

func (m *AssignTicketsRequest) GetAssignment() *Assignment {
	if m != nil {
		return m.Assignment
	}
	return nil
}

type AssignTicketsResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AssignTicketsResponse) Reset()         { *m = AssignTicketsResponse{} }
func (m *AssignTicketsResponse) String() string { return proto.CompactTextString(m) }
func (*AssignTicketsResponse) ProtoMessage()    {}
func (*AssignTicketsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8dab762378f455cd, []int{4}
}

func (m *AssignTicketsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AssignTicketsResponse.Unmarshal(m, b)
}
func (m *AssignTicketsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AssignTicketsResponse.Marshal(b, m, deterministic)
}
func (m *AssignTicketsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssignTicketsResponse.Merge(m, src)
}
func (m *AssignTicketsResponse) XXX_Size() int {
	return xxx_messageInfo_AssignTicketsResponse.Size(m)
}
func (m *AssignTicketsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AssignTicketsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AssignTicketsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("api.FunctionConfig_Type", FunctionConfig_Type_name, FunctionConfig_Type_value)
	proto.RegisterType((*FunctionConfig)(nil), "api.FunctionConfig")
	proto.RegisterType((*FetchMatchesRequest)(nil), "api.FetchMatchesRequest")
	proto.RegisterType((*FetchMatchesResponse)(nil), "api.FetchMatchesResponse")
	proto.RegisterType((*AssignTicketsRequest)(nil), "api.AssignTicketsRequest")
	proto.RegisterType((*AssignTicketsResponse)(nil), "api.AssignTicketsResponse")
}

func init() { proto.RegisterFile("api/backend.proto", fileDescriptor_8dab762378f455cd) }

var fileDescriptor_8dab762378f455cd = []byte{
	// 665 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x54, 0xd1, 0x52, 0xd3, 0x4c,
	0x18, 0xfd, 0x93, 0x96, 0x02, 0xcb, 0x2f, 0xc2, 0x82, 0x5a, 0xa2, 0xa3, 0x31, 0xea, 0x4c, 0xa7,
	0xd0, 0x2c, 0x44, 0x2e, 0x9c, 0x3a, 0xce, 0x50, 0x10, 0x1c, 0x66, 0x50, 0x99, 0xc0, 0x78, 0xe1,
	0x8d, 0x93, 0x26, 0x5f, 0x93, 0x95, 0x66, 0x77, 0xcd, 0x6e, 0x40, 0x6e, 0x79, 0x04, 0xbd, 0xf3,
	0x2d, 0x7c, 0x16, 0x6f, 0xbc, 0x97, 0x07, 0x71, 0xb2, 0xa1, 0xb4, 0x95, 0x5e, 0x75, 0x73, 0xbe,
	0xf3, 0x9d, 0x73, 0xf6, 0xdb, 0xdd, 0xa2, 0xc5, 0x40, 0x50, 0xd2, 0x0d, 0xc2, 0x13, 0x60, 0x91,
	0x2b, 0x32, 0xae, 0x38, 0xae, 0x04, 0x82, 0x5a, 0xb8, 0xc0, 0x53, 0x90, 0x32, 0x88, 0x41, 0x96,
	0x05, 0xeb, 0x41, 0xcc, 0x79, 0xdc, 0x07, 0x52, 0x94, 0x02, 0xc6, 0xb8, 0x0a, 0x14, 0xe5, 0x6c,
	0x50, 0x5d, 0xd3, 0x3f, 0x61, 0x2b, 0x06, 0xd6, 0x92, 0x67, 0x41, 0x1c, 0x43, 0x46, 0xb8, 0xd0,
	0x8c, 0x9b, 0x6c, 0xe7, 0xc2, 0x40, 0xf3, 0x7b, 0x39, 0x0b, 0x0b, 0x6c, 0x87, 0xb3, 0x1e, 0x8d,
	0x31, 0x46, 0xd5, 0x84, 0x4b, 0x55, 0x37, 0x6c, 0xa3, 0x31, 0xeb, 0xeb, 0x75, 0x81, 0x09, 0x9e,
	0xa9, 0xba, 0x69, 0x1b, 0x8d, 0x29, 0x5f, 0xaf, 0xf1, 0x1a, 0xaa, 0xaa, 0x73, 0x01, 0xf5, 0x8a,
	0x6d, 0x34, 0xe6, 0xbd, 0xba, 0x1b, 0x08, 0xea, 0x8e, 0x4b, 0xb9, 0xc7, 0xe7, 0x02, 0x7c, 0xcd,
	0x72, 0x2c, 0x54, 0x2d, 0xbe, 0xf0, 0x0c, 0xaa, 0xbe, 0xf1, 0x0f, 0x77, 0x16, 0xfe, 0x2b, 0x56,
	0xfe, 0xee, 0xd1, 0xf1, 0x82, 0xe1, 0x70, 0xb4, 0xb4, 0x07, 0x2a, 0x4c, 0xde, 0x06, 0x2a, 0x4c,
	0x40, 0xfa, 0xf0, 0x25, 0x07, 0xa9, 0xf0, 0x2a, 0xaa, 0x85, 0x5a, 0x47, 0x47, 0x99, 0xf3, 0x96,
	0x26, 0x58, 0xf8, 0x57, 0x14, 0xbc, 0x8a, 0xa6, 0x45, 0xc6, 0x7b, 0xb4, 0x0f, 0x75, 0xd3, 0xae,
	0x34, 0xe6, 0xbc, 0x45, 0xcd, 0xd6, 0x92, 0x87, 0x65, 0xc1, 0x1f, 0x30, 0x9c, 0x17, 0x68, 0x79,
	0xdc, 0x50, 0x0a, 0xce, 0x24, 0x60, 0x1b, 0x4d, 0xa5, 0x05, 0x54, 0x37, 0xb4, 0x04, 0x1a, 0x4a,
	0xf8, 0x65, 0xc1, 0x89, 0xd0, 0x72, 0x47, 0x4a, 0x1a, 0xb3, 0x63, 0x1a, 0x9e, 0x80, 0xba, 0xce,
	0x7a, 0x1f, 0xcd, 0x2a, 0x8d, 0x7c, 0xa2, 0x91, 0xee, 0x9e, 0xf5, 0x67, 0x4a, 0x60, 0x3f, 0xc2,
	0x04, 0xa1, 0x40, 0x37, 0xa5, 0xc0, 0xca, 0x19, 0xce, 0x79, 0xb7, 0xb5, 0x76, 0xe7, 0x1a, 0xf6,
	0x47, 0x28, 0xce, 0x3d, 0x74, 0xe7, 0x1f, 0x97, 0x32, 0xa0, 0xf7, 0xc7, 0x40, 0xd3, 0xdb, 0xe5,
	0x2d, 0xc1, 0x14, 0xfd, 0x3f, 0xba, 0x09, 0x7c, 0x75, 0x02, 0x37, 0x07, 0x69, 0xad, 0x4c, 0xa8,
	0x94, 0x82, 0xce, 0xd3, 0x8b, 0x5f, 0x97, 0xdf, 0xcd, 0x87, 0xce, 0x0a, 0x39, 0xdd, 0x18, 0xdc,
	0x3f, 0x92, 0x96, 0xa4, 0x76, 0xaf, 0xe8, 0x68, 0x1b, 0x4d, 0x9c, 0xa2, 0x5b, 0x63, 0x79, 0xf0,
	0xca, 0x48, 0xfa, 0xf1, 0x49, 0x58, 0xd6, 0xa4, 0xd2, 0x95, 0xdb, 0x33, 0xed, 0xf6, 0xc8, 0xb1,
	0x46, 0xdd, 0xca, 0x31, 0xc9, 0x76, 0xb9, 0xff, 0xb6, 0xd1, 0xdc, 0xbe, 0x34, 0xbf, 0x75, 0x7e,
	0x9b, 0xf8, 0xe7, 0x70, 0xb3, 0xce, 0x3e, 0x42, 0xef, 0x05, 0x30, 0x5b, 0xe7, 0xc7, 0x77, 0x13,
	0xa5, 0x84, 0x6c, 0x13, 0xc2, 0x05, 0xb0, 0x96, 0x8e, 0xeb, 0x46, 0x70, 0x6a, 0x3d, 0x19, 0x7e,
	0xb7, 0x22, 0x2a, 0xc3, 0x5c, 0xca, 0xad, 0xf2, 0xb1, 0xc4, 0x19, 0xcf, 0x85, 0x74, 0x43, 0x9e,
	0x36, 0x3f, 0x20, 0xdc, 0x11, 0x41, 0x98, 0x80, 0xed, 0xb9, 0xeb, 0xf6, 0x01, 0x0d, 0xa1, 0x38,
	0xfa, 0xad, 0x81, 0x64, 0x4c, 0x55, 0x92, 0x77, 0x0b, 0x26, 0x29, 0x5b, 0x7b, 0x3c, 0x8b, 0x83,
	0x14, 0xe4, 0x88, 0x19, 0xe9, 0xf6, 0x79, 0x97, 0xa4, 0x81, 0x54, 0x90, 0x91, 0x83, 0xfd, 0x9d,
	0xdd, 0x77, 0x47, 0xbb, 0x5e, 0x65, 0xc3, 0x5d, 0x6f, 0x9a, 0x86, 0xe9, 0x2d, 0x04, 0x42, 0xf4,
	0x69, 0xa8, 0xdf, 0x19, 0xf9, 0x2c, 0x39, 0x6b, 0xdf, 0x40, 0xfc, 0x97, 0xa8, 0xb2, 0xb9, 0xbe,
	0x89, 0x37, 0x51, 0xd3, 0x07, 0x95, 0x67, 0x0c, 0x22, 0xfb, 0x2c, 0x01, 0x66, 0xab, 0x04, 0xec,
	0x0c, 0x24, 0xcf, 0xb3, 0x10, 0xec, 0x88, 0x83, 0xb4, 0x19, 0x57, 0x36, 0x7c, 0xa5, 0x52, 0xb9,
	0xb8, 0x86, 0xaa, 0x3f, 0x4c, 0x63, 0x3a, 0x7b, 0x85, 0xea, 0xc3, 0x61, 0xd8, 0xaf, 0x79, 0x98,
	0x17, 0xb7, 0x46, 0xab, 0xe3, 0xc7, 0x93, 0x47, 0x43, 0x24, 0x55, 0x40, 0x22, 0x1e, 0x4a, 0xf2,
	0xb1, 0x26, 0x4e, 0x62, 0x22, 0xba, 0xdd, 0x9a, 0xfe, 0x0b, 0x78, 0xfe, 0x37, 0x00, 0x00, 0xff,
	0xff, 0xc4, 0x3a, 0xc0, 0xd7, 0x7c, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BackendClient is the client API for Backend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BackendClient interface {
	// FetchMatch triggers execution of the specfied MatchFunction for each of the
	// specified MatchProfiles. Each MatchFunction execution returns a set of
	// proposals which are then evaluated to generate results. FetchMatch method
	// streams these results back to the caller.
	FetchMatches(ctx context.Context, in *FetchMatchesRequest, opts ...grpc.CallOption) (*FetchMatchesResponse, error)
	// AssignTickets sets the specified Assignment on the Tickets for the Ticket
	// IDs passed.
	AssignTickets(ctx context.Context, in *AssignTicketsRequest, opts ...grpc.CallOption) (*AssignTicketsResponse, error)
}

type backendClient struct {
	cc *grpc.ClientConn
}

func NewBackendClient(cc *grpc.ClientConn) BackendClient {
	return &backendClient{cc}
}

func (c *backendClient) FetchMatches(ctx context.Context, in *FetchMatchesRequest, opts ...grpc.CallOption) (*FetchMatchesResponse, error) {
	out := new(FetchMatchesResponse)
	err := c.cc.Invoke(ctx, "/api.Backend/FetchMatches", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) AssignTickets(ctx context.Context, in *AssignTicketsRequest, opts ...grpc.CallOption) (*AssignTicketsResponse, error) {
	out := new(AssignTicketsResponse)
	err := c.cc.Invoke(ctx, "/api.Backend/AssignTickets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackendServer is the server API for Backend service.
type BackendServer interface {
	// FetchMatch triggers execution of the specfied MatchFunction for each of the
	// specified MatchProfiles. Each MatchFunction execution returns a set of
	// proposals which are then evaluated to generate results. FetchMatch method
	// streams these results back to the caller.
	FetchMatches(context.Context, *FetchMatchesRequest) (*FetchMatchesResponse, error)
	// AssignTickets sets the specified Assignment on the Tickets for the Ticket
	// IDs passed.
	AssignTickets(context.Context, *AssignTicketsRequest) (*AssignTicketsResponse, error)
}

func RegisterBackendServer(s *grpc.Server, srv BackendServer) {
	s.RegisterService(&_Backend_serviceDesc, srv)
}

func _Backend_FetchMatches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchMatchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).FetchMatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Backend/FetchMatches",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).FetchMatches(ctx, req.(*FetchMatchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_AssignTickets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignTicketsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).AssignTickets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Backend/AssignTickets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).AssignTickets(ctx, req.(*AssignTicketsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backend_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Backend",
	HandlerType: (*BackendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchMatches",
			Handler:    _Backend_FetchMatches_Handler,
		},
		{
			MethodName: "AssignTickets",
			Handler:    _Backend_AssignTickets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/backend.proto",
}
