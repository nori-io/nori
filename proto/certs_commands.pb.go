// Code generated by protoc-gen-go. DO NOT EDIT.
// source: certs_commands.proto

package commands

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// []bs = {len_hmac(1 byte), hmac..., body_file...}
type CertsUploadRequest struct {
	Pem                  []byte   `protobuf:"bytes,1,opt,name=pem,proto3" json:"pem,omitempty"`
	Key                  []byte   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CertsUploadRequest) Reset()         { *m = CertsUploadRequest{} }
func (m *CertsUploadRequest) String() string { return proto.CompactTextString(m) }
func (*CertsUploadRequest) ProtoMessage()    {}
func (*CertsUploadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8839626868052056, []int{0}
}

func (m *CertsUploadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CertsUploadRequest.Unmarshal(m, b)
}
func (m *CertsUploadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CertsUploadRequest.Marshal(b, m, deterministic)
}
func (m *CertsUploadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CertsUploadRequest.Merge(m, src)
}
func (m *CertsUploadRequest) XXX_Size() int {
	return xxx_messageInfo_CertsUploadRequest.Size(m)
}
func (m *CertsUploadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CertsUploadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CertsUploadRequest proto.InternalMessageInfo

func (m *CertsUploadRequest) GetPem() []byte {
	if m != nil {
		return m.Pem
	}
	return nil
}

func (m *CertsUploadRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func init() {
	proto.RegisterType((*CertsUploadRequest)(nil), "commands.CertsUploadRequest")
}

func init() { proto.RegisterFile("certs_commands.proto", fileDescriptor_8839626868052056) }

var fileDescriptor_8839626868052056 = []byte{
	// 102 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x49, 0x4e, 0x2d, 0x2a,
	0x29, 0x8e, 0x4f, 0xce, 0xcf, 0xcd, 0x4d, 0xcc, 0x4b, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x80, 0xf1, 0x95, 0x2c, 0xb8, 0x84, 0x9c, 0x41, 0x2a, 0x42, 0x0b, 0x72, 0xf2, 0x13,
	0x53, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x04, 0xb8, 0x98, 0x0b, 0x52, 0x73, 0x25,
	0x18, 0x15, 0x18, 0x35, 0x78, 0x82, 0x40, 0x4c, 0x90, 0x48, 0x76, 0x6a, 0xa5, 0x04, 0x13, 0x44,
	0x24, 0x3b, 0xb5, 0x32, 0x89, 0x0d, 0x6c, 0x94, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x71, 0x2a,
	0xdc, 0x7d, 0x62, 0x00, 0x00, 0x00,
}
