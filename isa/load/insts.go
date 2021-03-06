package load

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// Load instructions
	Lb  = isa.Op{Name: "lb",  F3:0b000, Functs: 0b0000000_000, Opcode: isa.Opcodes.L}
	Lbu = isa.Op{Name: "lbu", F3:0b100, Functs: 0b0000000_100, Opcode: isa.Opcodes.L}
	Lh  = isa.Op{Name: "lh",  F3:0b001, Functs: 0b0000000_001, Opcode: isa.Opcodes.L}
	Lhu = isa.Op{Name: "lhu", F3:0b101, Functs: 0b0000000_101, Opcode: isa.Opcodes.L}
	Lw  = isa.Op{Name: "lw",  F3:0b110, Functs: 0b0000000_110, Opcode: isa.Opcodes.L}
)
