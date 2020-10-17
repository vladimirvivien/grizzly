package ctrlunit

import (
	"log"

	"github.com/vladimirvivien/grizzly/isa"
)

// decodeR decodes register format instructions:
//
//   fn7     RD2   RD1   fn3 RS    OPCODE
//   0000000_00010_00001_000_00101_0110011
func decodeR(i isa.Inst) *isa.RFields {
	fields := isa.NewRFields()
	fields.Opcode = i & 0x7F
	fields.Rd = (i >> 7) & 0x1F
	fields.Funct3 = (i >> 12) & 0x7
	fields.Rs1 = (i >> 15) & 0x1F
	fields.Rs2 = (i >> 20) & 0x1F
	fields.Funct7 = (i >> 25) & 0x7F
	log.Printf(
		"decoder: decoded R fields {opcode=%07b, rd=%05b, funct3=%03b, rs1=%07b, rs2=%07b, funct7=%07b}",
		fields.Opcode, fields.Rd, fields.Funct3, fields.Rs1, fields.Rs2, fields.Funct7,
	)
	return fields
}

// decodeRI decodes register immediate format
//
// Decode with immediate value
//   imm[32:20]   RD1   fn3 RS    OPCODE
//   000000000010_00001_000_00101_0110011
//
// Decode with shift operations
//   fn7     Shift RD1   fn3 RS    OPCODE
//   0000000_00010_00001_000_00101_0110011
func decodeRI(i isa.Inst) *isa.RIFields {
	fields := isa.NewRIFields()
	fields.Opcode = i & 0x7F
	fields.Rd = (i >> 7) & 0x1F
	fields.Funct3 = (i >> 12) & 0x7
	fields.Rs1 = (i >> 15) & 0x1F
	switch fields.Funct3 {
	case 0b001, 0b101:
		fields.Shift = (i >> 20) & 0x1F
		fields.Funct7 = (i >> 25) & 0x7F
	default:
		fields.Imm = (i >> 20) & 0xFFF
	}
	log.Printf(
		"decoder: decoded I fields {opcode=%07b, rd=%05b, funct3=%03b, rs1=%07b, shift=%05b, funct7=%07b, imm=%010b}",
		fields.Opcode, fields.Rd, fields.Funct3, fields.Rs1, fields.Shift, fields.Funct7, fields.Imm,
	)
	return fields
}
