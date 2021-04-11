package jump

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// Load instructions
	Jal  = isa.Op{Name: "jal",  Opcode: isa.Opcodes.J}
	Jalr = isa.Op{Name: "lbu", F3:0, Opcode: isa.Opcodes.JI}
)

