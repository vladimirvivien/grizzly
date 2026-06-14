//go:build rv64 || rv64i

package branch

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

// Decode decodes 32-bit branch instruction format
func Decode(i datapath.XWord) datapath.OpFields {
	var fields datapath.OpFields
	fields.Opcode = uint8(i & 0x7F)
	imm11 := (i >> 7) & 0x1
	imm4_1 := (i >> 8) & 0xF
	fields.Funct3 = uint8((i >> 12) & 0x7)
	fields.Rs1 = uint8((i >> 15) & 0x1F)
	fields.Rs2 = uint8((i >> 20) & 0x1F)
	imm10_5 := (i >> 25) & 0x3F
	imm12 := (i >> 31) & 0x1
	fields.Imm = uint32(imm4_1 | imm10_5 << 4 | imm11 << 10 | imm12 << 11)
	return fields
}
