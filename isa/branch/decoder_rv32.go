//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package branch

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

// Decode decodes 32-bit branch instruction format
//
// 31_......24....19....14..11...7_......0
// I I[10:5] RS2   RS1  F3  Imm  I OPCODE
// 0_000000_00000_00000_000_0000_0_0000000
//
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
	// imm12 | imm11 | imm10_5 | imm4_1
	fields.Imm = imm4_1 | imm10_5 << 4 | imm11 << 10 | imm12 << 11
	return fields
}

