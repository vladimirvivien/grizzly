package branch

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// Branch instructions
	Beq  = isa.Op{Name: "beq", F3: 0, Opcode: isa.Opcodes.B}
	Bne  = isa.Op{Name: "bne", F3: 0b001, Opcode: isa.Opcodes.B}
	Blt  = isa.Op{Name: "blt", F3: 0b100, Opcode: isa.Opcodes.B}
	Bge  = isa.Op{Name: "bge", F3: 0b101, Opcode: isa.Opcodes.B}
	Bltu  = isa.Op{Name: "bltu", F3: 0b110, Opcode: isa.Opcodes.B}
	Bgeu = isa.Op{Name: "bgeu", F3: 0b111, Opcode: isa.Opcodes.B}
)
