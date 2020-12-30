package store

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// Load instructions
	Sb  = isa.Op{Name: "sb", Functs: uint32(0b0000000_000), Opcode: isa.Opcodes.S}
	Sh  = isa.Op{Name: "sh", Functs: uint32(0b0000000_001), Opcode: isa.Opcodes.S}
	Sw  = isa.Op{Name: "sw", Functs: uint32(0b0000000_010), Opcode: isa.Opcodes.S}
)
