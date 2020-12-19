package load

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// Load instructions
	Lb  = isa.Op{Name: "lb", Functs: uint32(0b0000000_000), Opcode: isa.Opcodes.L}
	Lbu = isa.Op{Name: "lbu", Functs: uint32(0b0000000_100), Opcode: isa.Opcodes.L}
	Lh  = isa.Op{Name: "lh", Functs: uint32(0b0000000_001), Opcode: isa.Opcodes.L}
	Lhu = isa.Op{Name: "lhu", Functs: uint32(0b0000000_101), Opcode: isa.Opcodes.L}
	Lw  = isa.Op{Name: "lw", Functs: uint32(0b0000000_110), Opcode: isa.Opcodes.L}
)
