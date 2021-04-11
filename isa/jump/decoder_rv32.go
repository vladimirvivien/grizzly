package jump

import (
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

// Decode decodes 32-bit jump (jal) instruction format
//
// 31................12......7__......0
//        Imm[19:0]      RD  OPCODE
// 00000000000000000000__00000__0000000
//
// Jump Immediate (jalr)
// 31........20__...15__.12__....7........0
//     Imm       RS1    fn3  RD     OPCODE
// 000000000010__00001__000__00101__0110011
func Decode(i datapath.XWord) datapath.OpFields {
	var fields datapath.OpFields
	fields.Opcode = uint8(i & 0x7F)
	switch fields.Opcode {
	case isa.Opcodes.J:
		fields.Rd = uint8((i >> 7) & 0x1F)
		fields.Imm = (i >> 12) & 0xFFFFF
	case isa.Opcodes.JI:
		fields.Rd = uint8((i >> 7) & 0x1F)
		fields.Funct3 = uint8((i >> 12) & 0x7)
		fields.Rs1 = uint8((i >> 15) & 0x1F)
		fields.Imm = (i >> 20) & 0xFFF
	}

	return fields
}
