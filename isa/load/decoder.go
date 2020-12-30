package load

import (
	"log"

	"github.com/vladimirvivien/grizzly/isa"
)

// LoadFields instruction fields for load instruction
type Fields struct {
	*isa.IntegerBaseFields
	Imm uint32
}

func newFields() *Fields {
	return &Fields{IntegerBaseFields: isa.NewIntegerBaseFields()}
}

// Decode decodes load instruction format
//
// 31............20.....15...12.....07.......0
//   imm[11:0]     RS1    fn3  RD     OPCODE
//   000000000000__00000__xxx__00000__0000011
//
func Decode(i isa.Inst) *Fields {
	fields := newFields()
	fields.Opcode = i & 0x7F
	fields.Rd = (i >> 7) & 0x1F
	fields.Funct3 = (i >> 12) & 0x7
	fields.Rs1 = (i >> 15) & 0x1F
	fields.Imm = (i >> 20) & 0xFFF
	log.Printf(
		"decoder: decoded L fields {opcode=%07b, rd=%05b, funct3=%03b, rs1=%07b, imm=%010b}",
		fields.Opcode, fields.Rd, fields.Funct3, fields.Rs1, fields.Imm,
	)
	return fields
}
