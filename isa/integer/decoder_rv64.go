//go:build rv64 || rv64i

package integer

import (
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

// Decode 64-bit integer format instructions
func Decode(i datapath.XWord) datapath.OpFields {
	var fields datapath.OpFields
	fields.Opcode = uint8(i & 0x7F)
	switch fields.Opcode {
	case isa.Opcodes.R:
		fields.Rd = uint8((i >> 7) & 0x1F)
		fields.Funct3 = uint8((i >> 12) & 0x7)
		fields.Rs1 = uint8((i >> 15) & 0x1F)
		fields.Rs2 = uint8((i >> 20) & 0x1F)
		fields.Funct7 = uint8((i >> 25) & 0x7F)
	case isa.Opcodes.RI:
		fields.Rd = uint8((i >> 7) & 0x1F)
		fields.Funct3 = uint8((i >> 12) & 0x7)
		fields.Rs1 = uint8((i >> 15) & 0x1F)
		switch fields.Funct3 {
		case Slli.F3, Srli.F3, Srai.F3:
			fields.Shift = uint8((i >> 20) & 0x3F)
			fields.Funct7 = uint8((i >> 26) & 0x3F)
		default:
			val := (i >> 20) & 0xFFF
			if (val & 0x800) != 0 {
				val |= 0xfffff000
			}
			fields.Imm = uint32(val)
		}
	}
	return fields
}
