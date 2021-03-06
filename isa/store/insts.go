package store

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// Load instructions
	Sb  = isa.Op{Name: "sb", F3:0b000, Functs: 0b0000000_000, Opcode: isa.Opcodes.S}
	Sh  = isa.Op{Name: "sh", F3:0b001, Functs: 0b0000000_001, Opcode: isa.Opcodes.S}
	Sw  = isa.Op{Name: "sw", F3:0b010, Functs: 0b0000000_010, Opcode: isa.Opcodes.S}
)
