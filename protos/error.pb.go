// Code generated by protoc-gen-go.
// source: error.proto
// DO NOT EDIT!

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	error.proto
	idp.proto
	supervisor.proto

It has these top-level messages:
	Error
	User
	Device
	AcquireCaptchaReq
	AcquireCaptchaRsp
	VerifyCaptchaReq
	VerifyCaptchaRsp
	RegisterUserReq
	RegisterUserRsp
	BindDeviceReq
	BindDeviceRsp
	VerifyDeviceReq
	VerifyDeviceRsp
	FarmerAccount
	FarmerOnLineReq
	FarmerOnLineRsp
	BlocksRange
	FarmerPingReq
	FarmerPingRsp
	FarmerConquerChallengeReq
	FarmerConquerChallengeRsp
	FarmerOffLineReq
	FarmerOffLineRsp
*/
package protos

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

type ErrorType int32

const (
	// everything is ok
	ErrorType_NONE_ERROR ErrorType = 0
	// request params is invalid
	ErrorType_INVALID_PARAM ErrorType = 1
	// system error
	ErrorType_INTERNAL_ERROR ErrorType = 2
	// user already sign up
	ErrorType_ALREADY_SIGNUP ErrorType = 3
	// captcha is invalid
	ErrorType_INVALID_CAPTCHA ErrorType = 4
	// user id is invalid
	ErrorType_INVALID_USERID ErrorType = 5
	// device can't be recognized
	ErrorType_INVALID_DEVICE ErrorType = 6
	// mac address already been taken by other device
	ErrorType_ALREADY_DEVICE_MAC ErrorType = 7
	// alias already benn taken by other device
	ErrorType_ALREADY_DEVICE_ALIAS ErrorType = 8
	// farmer online
	ErrorType_INVALID_STATE_FARMER_ONLINE ErrorType = 9
	// farmer offline
	ErrorType_INVALID_STATE_FARMER_OFFLINE ErrorType = 10
	// farmer challenge fail
	ErrorType_FARMER_CHALLENGE_FAIL ErrorType = 11
)

var ErrorType_name = map[int32]string{
	0:  "NONE_ERROR",
	1:  "INVALID_PARAM",
	2:  "INTERNAL_ERROR",
	3:  "ALREADY_SIGNUP",
	4:  "INVALID_CAPTCHA",
	5:  "INVALID_USERID",
	6:  "INVALID_DEVICE",
	7:  "ALREADY_DEVICE_MAC",
	8:  "ALREADY_DEVICE_ALIAS",
	9:  "INVALID_STATE_FARMER_ONLINE",
	10: "INVALID_STATE_FARMER_OFFLINE",
	11: "FARMER_CHALLENGE_FAIL",
}
var ErrorType_value = map[string]int32{
	"NONE_ERROR":                   0,
	"INVALID_PARAM":                1,
	"INTERNAL_ERROR":               2,
	"ALREADY_SIGNUP":               3,
	"INVALID_CAPTCHA":              4,
	"INVALID_USERID":               5,
	"INVALID_DEVICE":               6,
	"ALREADY_DEVICE_MAC":           7,
	"ALREADY_DEVICE_ALIAS":         8,
	"INVALID_STATE_FARMER_ONLINE":  9,
	"INVALID_STATE_FARMER_OFFLINE": 10,
	"FARMER_CHALLENGE_FAIL":        11,
}

func (x ErrorType) String() string {
	return proto.EnumName(ErrorType_name, int32(x))
}
func (ErrorType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Error struct {
	ErrorType ErrorType `protobuf:"varint,1,opt,name=errorType,enum=protos.ErrorType" json:"errorType,omitempty"`
	Message   string    `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Error)(nil), "protos.Error")
	proto.RegisterEnum("protos.ErrorType", ErrorType_name, ErrorType_value)
}

func init() { proto.RegisterFile("error.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 274 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x90, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x86, 0x6d, 0xb5, 0xad, 0x99, 0x62, 0xba, 0x1d, 0x3f, 0x88, 0x28, 0x58, 0xc4, 0x83, 0x78,
	0xe8, 0x41, 0xef, 0xc2, 0x90, 0x4c, 0xea, 0xc2, 0x66, 0x13, 0x36, 0x69, 0xc1, 0xd3, 0xa2, 0x10,
	0x3c, 0x49, 0x4a, 0xe2, 0xc5, 0xdf, 0xe3, 0x1f, 0x35, 0x1f, 0x0d, 0x8a, 0x78, 0x5a, 0x78, 0xde,
	0xe7, 0x7d, 0x59, 0x06, 0xa6, 0x79, 0x59, 0x16, 0xe5, 0x72, 0x5b, 0x16, 0x1f, 0x05, 0x8e, 0xdb,
	0xa7, 0xba, 0x7e, 0x84, 0x11, 0x37, 0x18, 0x6f, 0xc0, 0x69, 0xf3, 0xec, 0x73, 0x9b, 0x7b, 0x83,
	0xc5, 0xe0, 0xd6, 0xbd, 0x9f, 0x77, 0x6e, 0xb5, 0xe4, 0x3e, 0xc0, 0x19, 0x4c, 0xde, 0xf3, 0xaa,
	0x7a, 0x79, 0xcb, 0xbd, 0x61, 0xed, 0x38, 0x77, 0x5f, 0x43, 0x70, 0x7e, 0x62, 0x17, 0x40, 0xc7,
	0x9a, 0x2d, 0x1b, 0x13, 0x1b, 0xb1, 0x87, 0x73, 0x38, 0x92, 0x7a, 0x43, 0x4a, 0x06, 0x36, 0x21,
	0x43, 0x91, 0x18, 0x20, 0x82, 0x2b, 0x75, 0xc6, 0x46, 0x93, 0xda, 0x69, 0xc3, 0x86, 0x91, 0x32,
	0x4c, 0xc1, 0xb3, 0x4d, 0xe5, 0x4a, 0xaf, 0x13, 0xb1, 0x8f, 0xc7, 0x30, 0xeb, 0xab, 0x3e, 0x25,
	0x99, 0xff, 0x44, 0xe2, 0xa0, 0x2b, 0x77, 0x70, 0x9d, 0xb2, 0x91, 0x81, 0x18, 0xfd, 0x66, 0x01,
	0x6f, 0xa4, 0xcf, 0x62, 0x8c, 0x67, 0x80, 0xfd, 0x60, 0xc7, 0x6c, 0x44, 0xbe, 0x98, 0xa0, 0x07,
	0x27, 0x7f, 0x78, 0xdd, 0xa3, 0x54, 0x1c, 0xe2, 0x15, 0x5c, 0xf4, 0x2b, 0x69, 0x46, 0x19, 0xdb,
	0x90, 0x4c, 0xc4, 0xc6, 0xc6, 0x5a, 0x49, 0xcd, 0xc2, 0xc1, 0x05, 0x5c, 0xfe, 0x2f, 0x84, 0x61,
	0x6b, 0x00, 0x9e, 0xc3, 0xe9, 0x8e, 0xd5, 0x9f, 0x55, 0x8a, 0xf5, 0xaa, 0x91, 0xa4, 0x12, 0xd3,
	0xd7, 0xee, 0xda, 0x0f, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x2c, 0x3a, 0x7f, 0x10, 0x83, 0x01,
	0x00, 0x00,
}
