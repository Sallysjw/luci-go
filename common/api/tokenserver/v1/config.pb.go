// Code generated by protoc-gen-go.
// source: config.proto
// DO NOT EDIT!

package tokenserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// TokenServerConfig is read from tokenserver.cfg in luci-config.
type TokenServerConfig struct {
	// List of CAs we trust.
	CertificateAuthority []*CertificateAuthorityConfig `protobuf:"bytes,1,rep,name=certificate_authority" json:"certificate_authority,omitempty"`
}

func (m *TokenServerConfig) Reset()                    { *m = TokenServerConfig{} }
func (m *TokenServerConfig) String() string            { return proto.CompactTextString(m) }
func (*TokenServerConfig) ProtoMessage()               {}
func (*TokenServerConfig) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *TokenServerConfig) GetCertificateAuthority() []*CertificateAuthorityConfig {
	if m != nil {
		return m.CertificateAuthority
	}
	return nil
}

// CertificateAuthorityConfig defines a single CA we trust.
//
// Such CA issues certificates for nodes that use The Token Service. Each node
// has a private key and certificate with Common Name set to the FQDN of this
// node, e.g. "CN=slave43-c1.c.chromecompute.google.com.internal".
//
// The Token Server uses this CN to derive a name of a service account to
// associate with a node. It splits FQDN into a hostname ("slave43-c1") and
// a domain name ("c.chromecompute.google.com.internal"), searches for a domain
// name in "known_domains" map, and creates a service account in a Cloud Project
// specified there: <hostname>@<project-id>.iam.gserviceaccount.com.
//
// Note that we can't put FQDN in the service account email, since it is limited
// in length and doesn't allow '.' in it.
type CertificateAuthorityConfig struct {
	Cn       string `protobuf:"bytes,1,opt,name=cn" json:"cn,omitempty"`
	CertPath string `protobuf:"bytes,2,opt,name=cert_path" json:"cert_path,omitempty"`
	CrlUrl   string `protobuf:"bytes,3,opt,name=crl_url" json:"crl_url,omitempty"`
	UseOauth bool   `protobuf:"varint,4,opt,name=use_oauth" json:"use_oauth,omitempty"`
	// KnownDomains describes what cloud project to use for nodes in particular
	// domains.
	KnownDomains map[string]*DomainConfig `protobuf:"bytes,5,rep,name=known_domains" json:"known_domains,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *CertificateAuthorityConfig) Reset()                    { *m = CertificateAuthorityConfig{} }
func (m *CertificateAuthorityConfig) String() string            { return proto.CompactTextString(m) }
func (*CertificateAuthorityConfig) ProtoMessage()               {}
func (*CertificateAuthorityConfig) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *CertificateAuthorityConfig) GetKnownDomains() map[string]*DomainConfig {
	if m != nil {
		return m.KnownDomains
	}
	return nil
}

// DomainConfig is used inside CertificateAuthorityConfig.
type DomainConfig struct {
	// CloudProjectName is a name of Google Cloud Project to create service
	// accounts in.
	//
	// The Token Server's own service account must have Editor permission in this
	// project.
	CloudProjectName string `protobuf:"bytes,1,opt,name=cloud_project_name" json:"cloud_project_name,omitempty"`
}

func (m *DomainConfig) Reset()                    { *m = DomainConfig{} }
func (m *DomainConfig) String() string            { return proto.CompactTextString(m) }
func (*DomainConfig) ProtoMessage()               {}
func (*DomainConfig) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func init() {
	proto.RegisterType((*TokenServerConfig)(nil), "tokenserver.TokenServerConfig")
	proto.RegisterType((*CertificateAuthorityConfig)(nil), "tokenserver.CertificateAuthorityConfig")
	proto.RegisterType((*DomainConfig)(nil), "tokenserver.DomainConfig")
}

var fileDescriptor1 = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x91, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0x69, 0xeb, 0xd4, 0xbd, 0x4e, 0xa4, 0x01, 0xa1, 0xf6, 0x24, 0xbd, 0x58, 0x3c, 0xf4,
	0x30, 0x2f, 0xea, 0x4d, 0xa6, 0x5e, 0x04, 0x11, 0xf4, 0xe6, 0x21, 0xc4, 0xec, 0xcd, 0xc5, 0x76,
	0xc9, 0x48, 0x93, 0x49, 0x8f, 0xfe, 0xe7, 0xb6, 0x89, 0xc3, 0x89, 0x08, 0x3b, 0xbe, 0xef, 0x7d,
	0xdf, 0xf7, 0x7e, 0xf0, 0x60, 0xc4, 0x95, 0x9c, 0x89, 0xb7, 0x72, 0xa9, 0x95, 0x51, 0x24, 0x36,
	0xaa, 0x42, 0xd9, 0xa0, 0x5e, 0xa1, 0xce, 0x5f, 0x20, 0x79, 0xee, 0xc7, 0x27, 0x37, 0x4e, 0x9c,
	0x8f, 0xdc, 0xc1, 0x11, 0x47, 0x6d, 0xc4, 0x4c, 0x70, 0x66, 0x90, 0x32, 0x6b, 0xe6, 0x4a, 0x0b,
	0xd3, 0xa6, 0xc1, 0x49, 0x54, 0xc4, 0xe3, 0xd3, 0x72, 0xa3, 0xa1, 0x9c, 0xfc, 0x38, 0xaf, 0xd7,
	0x46, 0xdf, 0x93, 0x7f, 0x86, 0x90, 0xfd, 0xbf, 0x26, 0x00, 0x21, 0x97, 0x5d, 0x67, 0x50, 0x0c,
	0x49, 0x02, 0xc3, 0xfe, 0x24, 0x5d, 0x32, 0x33, 0x4f, 0x43, 0x27, 0x1d, 0xc2, 0x1e, 0xd7, 0x35,
	0xb5, 0xba, 0x4e, 0xa3, 0xb5, 0xc7, 0x36, 0x48, 0x55, 0xcf, 0x93, 0xee, 0x74, 0xd2, 0x3e, 0x79,
	0x84, 0x83, 0x4a, 0xaa, 0x0f, 0x49, 0xa7, 0x6a, 0xc1, 0x84, 0x6c, 0xd2, 0x81, 0x23, 0xbc, 0xdc,
	0x92, 0xb0, 0xbc, 0xef, 0xc3, 0x37, 0x3e, 0x7b, 0x2b, 0x8d, 0x6e, 0xb3, 0x07, 0x48, 0xfe, 0x88,
	0x24, 0x86, 0xa8, 0xc2, 0xf6, 0x1b, 0xb5, 0x80, 0xc1, 0x8a, 0xd5, 0x16, 0x1d, 0x66, 0x3c, 0x3e,
	0xfe, 0x75, 0xcb, 0xc7, 0x7c, 0xfb, 0x55, 0x78, 0x11, 0xe4, 0x67, 0x30, 0xda, 0xd4, 0x48, 0x06,
	0x84, 0xd7, 0xca, 0x4e, 0x69, 0xf7, 0x8c, 0x77, 0xe4, 0x86, 0x4a, 0xb6, 0x40, 0xdf, 0xfc, 0xba,
	0xeb, 0x1e, 0x74, 0xfe, 0x15, 0x00, 0x00, 0xff, 0xff, 0xae, 0x25, 0x23, 0xa8, 0xb0, 0x01, 0x00,
	0x00,
}
