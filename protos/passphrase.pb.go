// Code generated by protoc-gen-go.
// source: passphrase.proto
// DO NOT EDIT!

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PassphraseLanguage int32

const (
	PassphraseLanguage_English            PassphraseLanguage = 0
	PassphraseLanguage_SimplifiedChinese  PassphraseLanguage = 1
	PassphraseLanguage_TraditionalChinese PassphraseLanguage = 2
)

var PassphraseLanguage_name = map[int32]string{
	0: "English",
	1: "SimplifiedChinese",
	2: "TraditionalChinese",
}
var PassphraseLanguage_value = map[string]int32{
	"English":            0,
	"SimplifiedChinese":  1,
	"TraditionalChinese": 2,
}

func (x PassphraseLanguage) String() string {
	return proto.EnumName(PassphraseLanguage_name, int32(x))
}
func (PassphraseLanguage) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func init() {
	proto.RegisterEnum("protos.PassphraseLanguage", PassphraseLanguage_name, PassphraseLanguage_value)
}

func init() { proto.RegisterFile("passphrase.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 120 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x48, 0x2c, 0x2e,
	0x2e, 0xc8, 0x28, 0x4a, 0x2c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53,
	0xc5, 0x5a, 0x01, 0x5c, 0x42, 0x01, 0x70, 0x39, 0x9f, 0xc4, 0xbc, 0xf4, 0xd2, 0xc4, 0xf4, 0x54,
	0x21, 0x6e, 0x2e, 0x76, 0xd7, 0xbc, 0xf4, 0x9c, 0xcc, 0xe2, 0x0c, 0x01, 0x06, 0x21, 0x51, 0x2e,
	0xc1, 0xe0, 0xcc, 0xdc, 0x82, 0x9c, 0xcc, 0xb4, 0xcc, 0xd4, 0x14, 0xe7, 0x8c, 0xcc, 0xbc, 0xd4,
	0xe2, 0x54, 0x01, 0x46, 0x21, 0x31, 0x2e, 0xa1, 0x90, 0xa2, 0xc4, 0x94, 0xcc, 0x92, 0xcc, 0xfc,
	0xbc, 0xc4, 0x1c, 0x98, 0x38, 0x53, 0x12, 0xc4, 0x64, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xb1, 0x1f, 0xe0, 0x10, 0x74, 0x00, 0x00, 0x00,
}
