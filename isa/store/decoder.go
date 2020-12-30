package store

import (
	"log"

	"github.com/vladimirvivien/grizzly/isa"
)

// LoadFields instruction fields for load instruction
type Fields struct {
	*isa.IntegerBaseFields
	Rs2 uint32
	Imm uint32
	hiImm uint32
	loImm uint32
}

func newFields() *Fields {
	return &Fields{IntegerBaseFields: isa.NewIntegerBaseFields()}
}

// Decode decodes store instruction format
//
// 31.......25.....20.....15...12.....07.......0
//   I[11:5]  RS2    RS1    fn3  I[4:0] OPCODE
//   0000000__00000__00000__xxx__00000__0000011
//
// Where I are the 12-bit immediate value:
// hiImm = I[11:5] and loImm = I[4:0]
func Decode(i isa.Inst) *Fields {
	fields := newFields()
	fields.Opcode = i & 0x7F
	fields.loImm = (i >> 7) & 0x1F
	fields.Funct3 = (i >> 12) & 0x7
	fields.Rs1 = (i >> 15) & 0x1F
	fields.Rs2 = (i >> 20) & 0x1F
	fields.hiImm = (i >> 25) & 0xFFF
	fields.Imm = (fields.Imm | fields.hiImm) << 5
	fields.Imm = fields.Imm | fields.loImm

	log.Printf(
		"store: decoder: {opcode=%07b, loImm=%05b, funct3=%03b, rs1=%07b, rs2=%07b, hiImm=%010b (imm: %012b)}",
		fields.Opcode, fields.loImm, fields.Funct3, fields.Rs1, fields.Rs2, fields.hiImm, fields.Imm,
	)
	return fields
}

