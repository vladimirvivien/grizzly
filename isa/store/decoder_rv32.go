package store

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

// Decode decodes 32-bit store instruction format
//
// 31.......25.....20.....15...12.....07.......0
//   I[11:5]  RS2    RS1    fn3  I[4:0] OPCODE
//   0000000__00000__00000__xxx__00000__0000011
//
// Where I are the 12-bit immediate value:
// hiImm = I[11:5] and loImm = I[4:0]
func Decode(i datapath.XWord) datapath.OpFields {
	var fields datapath.OpFields
	fields.Opcode = uint8(i & 0x7F)
	loImm := (i >> 7) & 0x1F
	fields.Funct3 = uint8((i >> 12) & 0x7)
	fields.Rs1 = uint8((i >> 15) & 0x1F)
	fields.Rs2 = uint8((i >> 20) & 0x1F)
	hiImm := (i >> 25) & 0xFFF
	val := (hiImm << 5) | loImm
	if (val & 0x800) != 0 {
		val |= 0xfffff000
	}
	fields.Imm = val
	return fields
}

