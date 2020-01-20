// Code generated by protoc-gen-go. DO NOT EDIT.
// source: fileArtifactRepo.proto

package file_artifact_repo

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

type ArtifactCreated struct {
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Url                  string   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ArtifactCreated) Reset()         { *m = ArtifactCreated{} }
func (m *ArtifactCreated) String() string { return proto.CompactTextString(m) }
func (*ArtifactCreated) ProtoMessage()    {}
func (*ArtifactCreated) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8f44f7b76e30219, []int{0}
}

func (m *ArtifactCreated) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArtifactCreated.Unmarshal(m, b)
}
func (m *ArtifactCreated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArtifactCreated.Marshal(b, m, deterministic)
}
func (m *ArtifactCreated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArtifactCreated.Merge(m, src)
}
func (m *ArtifactCreated) XXX_Size() int {
	return xxx_messageInfo_ArtifactCreated.Size(m)
}
func (m *ArtifactCreated) XXX_DiscardUnknown() {
	xxx_messageInfo_ArtifactCreated.DiscardUnknown(m)
}

var xxx_messageInfo_ArtifactCreated proto.InternalMessageInfo

func (m *ArtifactCreated) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ArtifactCreated) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ArtifactCreated) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *ArtifactCreated) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type ArtifactDeprecated struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ArtifactDeprecated) Reset()         { *m = ArtifactDeprecated{} }
func (m *ArtifactDeprecated) String() string { return proto.CompactTextString(m) }
func (*ArtifactDeprecated) ProtoMessage()    {}
func (*ArtifactDeprecated) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8f44f7b76e30219, []int{1}
}

func (m *ArtifactDeprecated) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArtifactDeprecated.Unmarshal(m, b)
}
func (m *ArtifactDeprecated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArtifactDeprecated.Marshal(b, m, deterministic)
}
func (m *ArtifactDeprecated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArtifactDeprecated.Merge(m, src)
}
func (m *ArtifactDeprecated) XXX_Size() int {
	return xxx_messageInfo_ArtifactDeprecated.Size(m)
}
func (m *ArtifactDeprecated) XXX_DiscardUnknown() {
	xxx_messageInfo_ArtifactDeprecated.DiscardUnknown(m)
}

var xxx_messageInfo_ArtifactDeprecated proto.InternalMessageInfo

type ArtifactUnDeprecated struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ArtifactUnDeprecated) Reset()         { *m = ArtifactUnDeprecated{} }
func (m *ArtifactUnDeprecated) String() string { return proto.CompactTextString(m) }
func (*ArtifactUnDeprecated) ProtoMessage()    {}
func (*ArtifactUnDeprecated) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8f44f7b76e30219, []int{2}
}

func (m *ArtifactUnDeprecated) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArtifactUnDeprecated.Unmarshal(m, b)
}
func (m *ArtifactUnDeprecated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArtifactUnDeprecated.Marshal(b, m, deterministic)
}
func (m *ArtifactUnDeprecated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArtifactUnDeprecated.Merge(m, src)
}
func (m *ArtifactUnDeprecated) XXX_Size() int {
	return xxx_messageInfo_ArtifactUnDeprecated.Size(m)
}
func (m *ArtifactUnDeprecated) XXX_DiscardUnknown() {
	xxx_messageInfo_ArtifactUnDeprecated.DiscardUnknown(m)
}

var xxx_messageInfo_ArtifactUnDeprecated proto.InternalMessageInfo

type ArtifactDescriptionUpdated struct {
	Description          string   `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ArtifactDescriptionUpdated) Reset()         { *m = ArtifactDescriptionUpdated{} }
func (m *ArtifactDescriptionUpdated) String() string { return proto.CompactTextString(m) }
func (*ArtifactDescriptionUpdated) ProtoMessage()    {}
func (*ArtifactDescriptionUpdated) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8f44f7b76e30219, []int{3}
}

func (m *ArtifactDescriptionUpdated) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArtifactDescriptionUpdated.Unmarshal(m, b)
}
func (m *ArtifactDescriptionUpdated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArtifactDescriptionUpdated.Marshal(b, m, deterministic)
}
func (m *ArtifactDescriptionUpdated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArtifactDescriptionUpdated.Merge(m, src)
}
func (m *ArtifactDescriptionUpdated) XXX_Size() int {
	return xxx_messageInfo_ArtifactDescriptionUpdated.Size(m)
}
func (m *ArtifactDescriptionUpdated) XXX_DiscardUnknown() {
	xxx_messageInfo_ArtifactDescriptionUpdated.DiscardUnknown(m)
}

var xxx_messageInfo_ArtifactDescriptionUpdated proto.InternalMessageInfo

func (m *ArtifactDescriptionUpdated) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type CreateArtifact struct {
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Url                  string   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateArtifact) Reset()         { *m = CreateArtifact{} }
func (m *CreateArtifact) String() string { return proto.CompactTextString(m) }
func (*CreateArtifact) ProtoMessage()    {}
func (*CreateArtifact) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8f44f7b76e30219, []int{4}
}

func (m *CreateArtifact) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateArtifact.Unmarshal(m, b)
}
func (m *CreateArtifact) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateArtifact.Marshal(b, m, deterministic)
}
func (m *CreateArtifact) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateArtifact.Merge(m, src)
}
func (m *CreateArtifact) XXX_Size() int {
	return xxx_messageInfo_CreateArtifact.Size(m)
}
func (m *CreateArtifact) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateArtifact.DiscardUnknown(m)
}

var xxx_messageInfo_CreateArtifact proto.InternalMessageInfo

func (m *CreateArtifact) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *CreateArtifact) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateArtifact) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *CreateArtifact) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*ArtifactCreated)(nil), "file.artifact.repo.ArtifactCreated")
	proto.RegisterType((*ArtifactDeprecated)(nil), "file.artifact.repo.ArtifactDeprecated")
	proto.RegisterType((*ArtifactUnDeprecated)(nil), "file.artifact.repo.ArtifactUnDeprecated")
	proto.RegisterType((*ArtifactDescriptionUpdated)(nil), "file.artifact.repo.ArtifactDescriptionUpdated")
	proto.RegisterType((*CreateArtifact)(nil), "file.artifact.repo.CreateArtifact")
}

func init() { proto.RegisterFile("fileArtifactRepo.proto", fileDescriptor_d8f44f7b76e30219) }

var fileDescriptor_d8f44f7b76e30219 = []byte{
	// 194 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x90, 0xc1, 0xca, 0xc2, 0x30,
	0x0c, 0xc7, 0xd9, 0xb7, 0xf1, 0x89, 0x11, 0x54, 0xc2, 0x18, 0xc5, 0xd3, 0xd8, 0xc9, 0xd3, 0x2e,
	0xde, 0x05, 0xd1, 0x27, 0x18, 0xec, 0x01, 0xea, 0x96, 0x41, 0x61, 0xae, 0x25, 0xab, 0x3e, 0xbf,
	0xb4, 0x5a, 0x1d, 0x7a, 0xf6, 0x96, 0xfe, 0x12, 0xfe, 0xbf, 0x26, 0x90, 0x75, 0xaa, 0xa7, 0x03,
	0x5b, 0xd5, 0xc9, 0xc6, 0x56, 0x64, 0x74, 0x69, 0x58, 0x5b, 0x8d, 0xe8, 0x78, 0x29, 0x9f, 0x8d,
	0x92, 0xc9, 0xe8, 0x62, 0x84, 0x55, 0x98, 0x3c, 0x32, 0x49, 0x4b, 0x2d, 0x0a, 0x98, 0xdd, 0x88,
	0x47, 0xa5, 0x07, 0xf1, 0x97, 0x47, 0xdb, 0x79, 0x15, 0x9e, 0x88, 0x90, 0x0c, 0xf2, 0x42, 0x22,
	0xf2, 0xd8, 0xd7, 0xb8, 0x86, 0xf8, 0xca, 0xbd, 0x88, 0x3d, 0x72, 0x25, 0xe6, 0xb0, 0x68, 0x69,
	0x6c, 0x58, 0x19, 0xeb, 0x32, 0x12, 0xdf, 0x99, 0xa2, 0x22, 0x05, 0x0c, 0xd2, 0x13, 0x19, 0xa6,
	0xc6, 0x79, 0x8b, 0x0c, 0xd2, 0x40, 0xeb, 0x61, 0xc2, 0xf7, 0xb0, 0x79, 0x4f, 0xbf, 0x42, 0x6a,
	0xd3, 0xfa, 0xdf, 0x7e, 0xd8, 0xa2, 0x6f, 0x1b, 0xc3, 0xf2, 0xb1, 0x5a, 0x48, 0xf9, 0xfd, 0x86,
	0xe7, 0x7f, 0x7f, 0xf1, 0xdd, 0x3d, 0x00, 0x00, 0xff, 0xff, 0x8b, 0x15, 0x41, 0x08, 0x8b, 0x01,
	0x00, 0x00,
}
