package load

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

// Decode 32-bit load instruction format
//
// 31............20.....15...12.....07.......0
//   imm[11:0]     RS1    fn3  RD     OPCODE
//   000000000000__00000__xxx__00000__0000011
//
func Decode(i datapath.XWord) datapath.OpFields {
	var fields datapath.OpFields
	fields.Opcode = uint8(i & 0x7F)
	fields.Rd = uint8((i >> 7) & 0x1F)
	fields.Funct3 = uint8((i >> 12) & 0x7)
	fields.Rs1 = uint8((i >> 15) & 0x1F)
	fields.Imm = (i >> 20) & 0xFFF
	return fields
}
