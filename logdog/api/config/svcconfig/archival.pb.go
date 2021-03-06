// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/logdog/api/config/svcconfig/archival.proto
// DO NOT EDIT!

/*
Package svcconfig is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/logdog/api/config/svcconfig/archival.proto
	github.com/luci/luci-go/logdog/api/config/svcconfig/config.proto
	github.com/luci/luci-go/logdog/api/config/svcconfig/project.proto
	github.com/luci/luci-go/logdog/api/config/svcconfig/storage.proto
	github.com/luci/luci-go/logdog/api/config/svcconfig/transport.proto

It has these top-level messages:
	ArchiveIndexConfig
	Config
	Coordinator
	Collector
	Archivist
	ProjectConfig
	Storage
	Transport
*/
package svcconfig

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// ArchiveIndexConfig specifies how archive indexes should be generated.
//
// By default, each log entry will be present in the index. This is generally
// overkill; instead, the index can be more sparse at the expense of a slightly
// higher data load.
type ArchiveIndexConfig struct {
	// If not zero, the maximum number of stream indices between index entries.
	StreamRange int32 `protobuf:"varint,1,opt,name=stream_range,json=streamRange" json:"stream_range,omitempty"`
	// If not zero, the maximum number of prefix indices between index entries.
	PrefixRange int32 `protobuf:"varint,2,opt,name=prefix_range,json=prefixRange" json:"prefix_range,omitempty"`
	// If not zero, the maximum number of log data bytes between index entries.
	ByteRange int32 `protobuf:"varint,3,opt,name=byte_range,json=byteRange" json:"byte_range,omitempty"`
}

func (m *ArchiveIndexConfig) Reset()                    { *m = ArchiveIndexConfig{} }
func (m *ArchiveIndexConfig) String() string            { return proto.CompactTextString(m) }
func (*ArchiveIndexConfig) ProtoMessage()               {}
func (*ArchiveIndexConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*ArchiveIndexConfig)(nil), "svcconfig.ArchiveIndexConfig")
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/logdog/api/config/svcconfig/archival.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x72, 0x4a, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0x29, 0x4d, 0xce, 0x04, 0x13, 0xba, 0xe9, 0xf9,
	0xfa, 0x39, 0xf9, 0xe9, 0x29, 0xf9, 0xe9, 0xfa, 0x89, 0x05, 0x99, 0xfa, 0xc9, 0xf9, 0x79, 0x69,
	0x99, 0xe9, 0xfa, 0xc5, 0x65, 0xc9, 0x50, 0x56, 0x62, 0x51, 0x72, 0x46, 0x66, 0x59, 0x62, 0x8e,
	0x5e, 0x41, 0x51, 0x7e, 0x49, 0xbe, 0x10, 0x27, 0x5c, 0x46, 0xa9, 0x92, 0x4b, 0xc8, 0x11, 0x2c,
	0x99, 0xea, 0x99, 0x97, 0x92, 0x5a, 0xe1, 0x0c, 0x16, 0x15, 0x52, 0xe4, 0xe2, 0x29, 0x2e, 0x29,
	0x4a, 0x4d, 0xcc, 0x8d, 0x2f, 0x4a, 0xcc, 0x4b, 0x4f, 0x95, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x0d,
	0xe2, 0x86, 0x88, 0x05, 0x81, 0x84, 0x40, 0x4a, 0x0a, 0x8a, 0x52, 0xd3, 0x32, 0x2b, 0xa0, 0x4a,
	0x98, 0x20, 0x4a, 0x20, 0x62, 0x10, 0x25, 0xb2, 0x5c, 0x5c, 0x49, 0x95, 0x25, 0xa9, 0x50, 0x05,
	0xcc, 0x60, 0x05, 0x9c, 0x20, 0x11, 0xb0, 0x74, 0x12, 0x1b, 0xd8, 0x31, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x06, 0x73, 0xbd, 0x23, 0xd2, 0x00, 0x00, 0x00,
}
