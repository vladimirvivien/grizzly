package integer

import (
	"log"

	"github.com/vladimirvivien/grizzly/isa"
)
// Fields for integer instructions
type Fields struct {
	*isa.IntegerBaseFields
	Rs2    uint32 // register 2
	Funct7 uint32 // 7-bit funct field
}
// newFields returns *Fields
func newFields() *Fields {
	return &Fields{IntegerBaseFields: isa.NewIntegerBaseFields()}
}
func (f *Fields) Functs() uint32 {
	return Functs(f.Funct7, f.Funct3)
}

// Decode decodes integer format instructions:
//
//   fn7     RS2   RS1   fn3 RD    OPCODE
//   0000000_00010_00001_000_00101_0110011
func Decode(i isa.Inst) *Fields {
	fields := newFields()
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


