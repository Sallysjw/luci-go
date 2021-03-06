// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/tokenserver/api/token_file.proto
// DO NOT EDIT!

package tokenserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// TokenFile is representation of a token file on disk (serialized as JSON).
//
// The token file is consumed by whoever wishes to use machine tokens. It is
// intentionally made as simple as possible (e.g. uses unix timestamps instead
// of fancy protobuf ones).
type TokenFile struct {
	// Google OAuth2 access token of a machine service account.
	AccessToken string `protobuf:"bytes,1,opt,name=access_token" json:"access_token,omitempty"`
	// OAuth2 access token type, usually "Bearer".
	TokenType string `protobuf:"bytes,2,opt,name=token_type" json:"token_type,omitempty"`
	// Machine token understood by LUCI backends (alternative to access_token).
	LuciMachineToken string `protobuf:"bytes,3,opt,name=luci_machine_token" json:"luci_machine_token,omitempty"`
	// Unix timestamp (in seconds) when this token expires.
	//
	// The token file is expected to be updated before the token expires, see
	// 'next_update' for next expected update time.
	Expiry int64 `protobuf:"varint,4,opt,name=expiry" json:"expiry,omitempty"`
	// Unix timestamp of when this file was updated the last time.
	LastUpdate int64 `protobuf:"varint,5,opt,name=last_update" json:"last_update,omitempty"`
	// Unix timestamp of when this file is expected to be updated next time.
	NextUpdate int64 `protobuf:"varint,6,opt,name=next_update" json:"next_update,omitempty"`
	// Email of the associated service account.
	ServiceAccountEmail string `protobuf:"bytes,7,opt,name=service_account_email" json:"service_account_email,omitempty"`
	// Unique stable ID of the associated service account.
	ServiceAccountUniqueId string `protobuf:"bytes,8,opt,name=service_account_unique_id" json:"service_account_unique_id,omitempty"`
	// Any information tokend daemon wishes to associate with the token.
	//
	// Consumers of the token file should ignore this field. It is used
	// exclusively by tokend daemon.
	TokendState []byte `protobuf:"bytes,50,opt,name=tokend_state,proto3" json:"tokend_state,omitempty"`
}

func (m *TokenFile) Reset()                    { *m = TokenFile{} }
func (m *TokenFile) String() string            { return proto.CompactTextString(m) }
func (*TokenFile) ProtoMessage()               {}
func (*TokenFile) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func init() {
	proto.RegisterType((*TokenFile)(nil), "tokenserver.TokenFile")
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/tokenserver/api/token_file.proto", fileDescriptor2)
}

var fileDescriptor2 = []byte{
	// 262 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x91, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0xd5, 0x16, 0x02, 0xbd, 0x76, 0xb2, 0x04, 0x32, 0x0b, 0xaa, 0x3a, 0xb1, 0x90, 0x48,
	0xc0, 0xc0, 0xc0, 0xcc, 0x0f, 0x88, 0xd8, 0x2d, 0xd7, 0x39, 0xda, 0x13, 0x89, 0x6d, 0x62, 0x1b,
	0xb5, 0xbf, 0x9b, 0x3f, 0x40, 0xec, 0x54, 0x28, 0x45, 0x61, 0x89, 0x74, 0xef, 0x7d, 0x77, 0x7a,
	0x2f, 0x86, 0xe7, 0x2d, 0xf9, 0x5d, 0xd8, 0xe4, 0xca, 0x34, 0x45, 0x1d, 0x14, 0xa5, 0xcf, 0xfd,
	0xd6, 0x14, 0xde, 0x7c, 0xa0, 0x76, 0xd8, 0x7e, 0x61, 0x5b, 0x48, 0x4b, 0xfd, 0x2c, 0xde, 0xa9,
	0xc6, 0xdc, 0xb6, 0xc6, 0x1b, 0xb6, 0x18, 0x10, 0xeb, 0xef, 0x29, 0xcc, 0xdf, 0xe2, 0xfc, 0xda,
	0x01, 0x6c, 0x0d, 0x4b, 0xa9, 0x14, 0x3a, 0x27, 0x12, 0xc3, 0x27, 0xab, 0xc9, 0xdd, 0xbc, 0x3c,
	0xd1, 0xd8, 0x2d, 0x40, 0x7f, 0xd2, 0x1f, 0x2c, 0xf2, 0x69, 0x22, 0x06, 0x0a, 0xcb, 0x81, 0xc5,
	0x28, 0xa2, 0x91, 0x6a, 0x47, 0x1a, 0x8f, 0x97, 0x66, 0x89, 0x1b, 0x71, 0xd8, 0x35, 0x64, 0xb8,
	0xb7, 0xd4, 0x1e, 0xf8, 0x59, 0xc7, 0xcc, 0xca, 0xe3, 0xc4, 0x56, 0xb0, 0xa8, 0xa5, 0xf3, 0x22,
	0xd8, 0x4a, 0x7a, 0xe4, 0xe7, 0xc9, 0x1c, 0x4a, 0x91, 0xd0, 0xb8, 0xff, 0x25, 0xb2, 0x9e, 0x18,
	0x48, 0xec, 0x09, 0xae, 0x62, 0x4f, 0x52, 0x28, 0xba, 0x0e, 0x26, 0x68, 0x2f, 0xb0, 0x91, 0x54,
	0xf3, 0x8b, 0x14, 0x67, 0xdc, 0x64, 0x2f, 0x70, 0xf3, 0xd7, 0x08, 0x9a, 0x3e, 0x03, 0x0a, 0xaa,
	0xf8, 0x65, 0xda, 0xfc, 0x1f, 0x88, 0xff, 0x30, 0x15, 0xab, 0x84, 0xf3, 0x31, 0xd6, 0x43, 0xb7,
	0xb0, 0x2c, 0x4f, 0xb4, 0x4d, 0x96, 0x5e, 0xe2, 0xf1, 0x27, 0x00, 0x00, 0xff, 0xff, 0xf3, 0xf5,
	0xfb, 0x8f, 0xc5, 0x01, 0x00, 0x00,
}
