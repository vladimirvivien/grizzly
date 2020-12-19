package integer

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// RISC-32 R
	Add  = isa.Op{Name: "add", Functs: uint32(0b0000000_000), Opcode: isa.Opcodes.R}
	Sub  = isa.Op{Name: "sub", Functs: uint32(0b0100000_000), Opcode: isa.Opcodes.R}
	Sll  = isa.Op{Name: "sll", Functs: uint32(0b0000000_001), Opcode: isa.Opcodes.R}
	Slt  = isa.Op{Name: "slt", Functs: uint32(0b0000000_010), Opcode: isa.Opcodes.R}
	Sltu = isa.Op{Name: "sltu", Functs: uint32(0b0000000_011), Opcode: isa.Opcodes.R}
	Xor  = isa.Op{Name: "xor", Functs: uint32(0b0000000_100), Opcode: isa.Opcodes.R}
	Srl  = isa.Op{Name: "srl", Functs: uint32(0b0000000_101), Opcode: isa.Opcodes.R}
	Sra  = isa.Op{Name: "sra", Functs: uint32(0b0100000_101), Opcode: isa.Opcodes.R}
	Or   = isa.Op{Name: "or", Functs: uint32(0b0000000_110), Opcode: isa.Opcodes.R}
	And  = isa.Op{Name: "and", Functs: uint32(0b0000000_111), Opcode: isa.Opcodes.R}

	// RISC-32 Register-Immediate (RI)
	Addi  = isa.Op{Name: "addi", Functs: uint32(0b0000000_000), Opcode: isa.Opcodes.RI}
	Slli  = isa.Op{Name: "slli", Functs: uint32(0b0000000_001), Opcode: isa.Opcodes.RI}
	Slti  = isa.Op{Name: "slti", Functs: uint32(0b0000000_010), Opcode: isa.Opcodes.RI}
	Sltiu = isa.Op{Name: "sltiu", Functs: uint32(0b0000000_011), Opcode: isa.Opcodes.RI}
	Xori  = isa.Op{Name: "xori", Functs: uint32(0b0000000_100), Opcode: isa.Opcodes.RI}
	Srli  = isa.Op{Name: "srli", Functs: uint32(0b0000000_101), Opcode: isa.Opcodes.RI}
	Srai  = isa.Op{Name: "srai", Functs: uint32(0b0100000_101), Opcode: isa.Opcodes.RI}
	Ori   = isa.Op{Name: "ori", Functs: uint32(0b0000000_110), Opcode: isa.Opcodes.RI}
	Andi  = isa.Op{Name: "andi", Functs: uint32(0b0000000_111), Opcode: isa.Opcodes.RI}

	// RISC-32 M
	Mul    = isa.Op{Name: "mul", Functs: uint32(0b0000001_000), Opcode: isa.Opcodes.R}
	Mulh   = isa.Op{Name: "mulh", Functs: uint32(0b0000001_001), Opcode: isa.Opcodes.R}
	Mulhsu = isa.Op{Name: "mulhsu", Functs: uint32(0b0000001_010), Opcode: isa.Opcodes.R}
	Mulhu  = isa.Op{Name: "mulhu", Functs: uint32(0b0000001_011), Opcode: isa.Opcodes.R}
	Div    = isa.Op{Name: "div", Functs: uint32(0b0000001_100), Opcode: isa.Opcodes.R}
	Divu   = isa.Op{Name: "divu", Functs: uint32(0b0000001_101), Opcode: isa.Opcodes.R}
	Rem    = isa.Op{Name: "rem", Functs: uint32(0b0000001_110), Opcode: isa.Opcodes.R}
	Remu   = isa.Op{Name: "remu", Functs: uint32(0b0000001_111), Opcode: isa.Opcodes.R}
)
