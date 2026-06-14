//go:build rv64 || rv64i

package jump

import (
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

// Decode decodes 32-bit jump (jal/jalr) instruction format
func Decode(i datapath.XWord) datapath.OpFields {
	var fields datapath.OpFields
	fields.Opcode = uint8(i & 0x7F)
	switch fields.Opcode {
	case isa.Opcodes.J:
		fields.Rd = uint8((i >> 7) & 0x1F)
		imm20 := (i >> 31) & 0x1
		imm10_1 := (i >> 21) & 0x3FF
		imm11 := (i >> 20) & 0x1
		imm19_12 := (i >> 12) & 0xFF
		val := imm10_1 | imm11 << 10 | imm19_12 << 11 | imm20 << 19
		offset := val << 1
		if (offset & 0x100000) != 0 {
			offset |= 0xffe00000
		}
		fields.Imm = uint32(offset)
	case isa.Opcodes.JI:
		fields.Rd = uint8((i >> 7) & 0x1F)
		fields.Funct3 = uint8((i >> 12) & 0x7)
		fields.Rs1 = uint8((i >> 15) & 0x1F)
		val := (i >> 20) & 0xFFF
		if (val & 0x800) != 0 {
			val |= 0xfffff000
		}
		fields.Imm = uint32(val)
	}

	return fields
}
