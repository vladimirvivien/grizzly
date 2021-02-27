package integer

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// RISC-32 R
	Add  = isa.Op{Name: "add", Functs: 0b0000000_000, Opcode: isa.Opcodes.R}
	Sub  = isa.Op{Name: "sub", Functs: 0b0100000_000, Opcode: isa.Opcodes.R}
	Sll  = isa.Op{Name: "sll", Functs: 0b0000000_001, Opcode: isa.Opcodes.R}
	Slt  = isa.Op{Name: "slt", Functs: 0b0000000_010, Opcode: isa.Opcodes.R}
	Sltu = isa.Op{Name: "sltu", Functs: 0b0000000_011, Opcode: isa.Opcodes.R}
	Xor  = isa.Op{Name: "xor", Functs: 0b0000000_100, Opcode: isa.Opcodes.R}
	Srl  = isa.Op{Name: "srl", Functs: 0b0000000_101, Opcode: isa.Opcodes.R}
	Sra  = isa.Op{Name: "sra", Functs: 0b0100000_101, Opcode: isa.Opcodes.R}
	Or   = isa.Op{Name: "or", Functs: 0b0000000_110, Opcode: isa.Opcodes.R}
	And  = isa.Op{Name: "and", Functs: 0b0000000_111, Opcode: isa.Opcodes.R}

	// RISC-32 Register-Immediate (RI)
	Addi  = isa.Op{Name: "addi", Functs: 0b0000000_000, Opcode: isa.Opcodes.RI}
	Slli  = isa.Op{Name: "slli", Functs: 0b0000000_001, Opcode: isa.Opcodes.RI}
	Slti  = isa.Op{Name: "slti", Functs: 0b0000000_010, Opcode: isa.Opcodes.RI}
	Sltiu = isa.Op{Name: "sltiu", Functs: 0b0000000_011, Opcode: isa.Opcodes.RI}
	Xori  = isa.Op{Name: "xori", Functs: 0b0000000_100, Opcode: isa.Opcodes.RI}
	Srli  = isa.Op{Name: "srli", Functs: 0b0000000_101, Opcode: isa.Opcodes.RI}
	Srai  = isa.Op{Name: "srai", Functs: 0b0100000_101, Opcode: isa.Opcodes.RI}
	Ori   = isa.Op{Name: "ori", Functs: 0b0000000_110, Opcode: isa.Opcodes.RI}
	Andi  = isa.Op{Name: "andi", Functs: 0b0000000_111, Opcode: isa.Opcodes.RI}

	// RISC-32 M
	Mul    = isa.Op{Name: "mul", Functs: 0b0000001_000, Opcode: isa.Opcodes.R}
	Mulh   = isa.Op{Name: "mulh", Functs: 0b0000001_001, Opcode: isa.Opcodes.R}
	Mulhsu = isa.Op{Name: "mulhsu", Functs: 0b0000001_010, Opcode: isa.Opcodes.R}
	Mulhu  = isa.Op{Name: "mulhu", Functs: 0b0000001_011, Opcode: isa.Opcodes.R}
	Div    = isa.Op{Name: "div", Functs: 0b0000001_100, Opcode: isa.Opcodes.R}
	Divu   = isa.Op{Name: "divu", Functs: 0b0000001_101, Opcode: isa.Opcodes.R}
	Rem    = isa.Op{Name: "rem", Functs: 0b0000001_110, Opcode: isa.Opcodes.R}
	Remu   = isa.Op{Name: "remu", Functs: 0b0000001_111, Opcode: isa.Opcodes.R}
)
