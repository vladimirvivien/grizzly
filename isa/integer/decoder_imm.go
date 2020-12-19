package integer

import (
	"log"

	"github.com/vladimirvivien/grizzly/isa"
)

// FieldsImm for integer immediate ops
type FieldsImm struct {
	*isa.IntegerBaseFields
	Imm uint32
	// Shift when Funct3 = 001, 101
	Shift  uint32
	Funct7 uint32
}

// newFieldsImm returns *FieldsImm
func newFieldsImm() *FieldsImm {
	return &FieldsImm{IntegerBaseFields: isa.NewIntegerBaseFields()}
}

func (f *FieldsImm) Functs() uint32 {
	return Functs(f.Funct7, f.Funct3)
}

// DecodeImm decodes integer immediate instructions
//
// Ins with immediate value
//   imm[32:20]   RD1   fn3 RS    OPCODE
//   000000000010_00001_000_00101_0110011
//
// Inst with shift operations
//   fn7     Shift RD1   fn3 RS    OPCODE
//   0000000_00010_00001_000_00101_0110011
func DecodeImm(i isa.Inst) *FieldsImm {
	fields := newFieldsImm()
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
