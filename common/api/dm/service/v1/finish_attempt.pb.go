// Code generated by protoc-gen-go.
// source: finish_attempt.proto
// DO NOT EDIT!

package dm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// FinishAttemptReq sets the final result of an Attempt.
type FinishAttemptReq struct {
	// required
	Auth *Execution_Auth `protobuf:"bytes,1,opt,name=auth" json:"auth,omitempty"`
	// The result data for this Attempt. The `size` field is recalculated after
	// the data field is normalized, and may be omitted.
	Data *JsonResult `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *FinishAttemptReq) Reset()                    { *m = FinishAttemptReq{} }
func (m *FinishAttemptReq) String() string            { return proto.CompactTextString(m) }
func (*FinishAttemptReq) ProtoMessage()               {}
func (*FinishAttemptReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *FinishAttemptReq) GetAuth() *Execution_Auth {
	if m != nil {
		return m.Auth
	}
	return nil
}

func (m *FinishAttemptReq) GetData() *JsonResult {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*FinishAttemptReq)(nil), "dm.FinishAttemptReq")
}

func init() { proto.RegisterFile("finish_attempt.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 144 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x49, 0xcb, 0xcc, 0xcb,
	0x2c, 0xce, 0x88, 0x4f, 0x2c, 0x29, 0x49, 0xcd, 0x2d, 0x28, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x4a, 0xc9, 0x95, 0x12, 0x48, 0x2f, 0x4a, 0x2c, 0xc8, 0x88, 0x4f, 0x49, 0x2c, 0x49,
	0x84, 0x88, 0x2a, 0xc5, 0x71, 0x09, 0xb8, 0x81, 0x55, 0x3b, 0x42, 0x14, 0x07, 0xa5, 0x16, 0x0a,
	0xa9, 0x71, 0xb1, 0x24, 0x96, 0x96, 0x64, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x09, 0xe9,
	0xa5, 0xe4, 0xea, 0xb9, 0x56, 0xa4, 0x26, 0x97, 0x96, 0x64, 0xe6, 0xe7, 0xe9, 0x39, 0x02, 0x65,
	0x82, 0xc0, 0xf2, 0x42, 0x4a, 0x5c, 0x2c, 0x20, 0x93, 0x24, 0x98, 0xc0, 0xea, 0xf8, 0x40, 0xea,
	0xbc, 0x8a, 0xf3, 0xf3, 0x82, 0x52, 0x8b, 0x4b, 0x73, 0x4a, 0x82, 0xc0, 0x72, 0x49, 0x6c, 0x60,
	0x6b, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x5a, 0x1a, 0x51, 0x1c, 0x94, 0x00, 0x00, 0x00,
}
