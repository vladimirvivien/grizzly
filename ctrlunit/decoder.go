package ctrlunit

import (
	"github.com/vladimirvivien/grizzly/isa"
)

func decodeR(i isa.Inst) *isa.RFields {
	return &isa.RFields{
		Opcode: isa.Opcodes.R,
		Rd:     (i >> 7) & 0x1F,
		Funct3: (i >> 12) & 0x7,
		Rs1:    (i >> 15) & 0x1F,
		Rs2:    (i >> 20) & 0x1F,
		Funct7: (i >> 25) & 0x7F,
	}
}
