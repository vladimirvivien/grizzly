package integer

import (
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	// RISC-32 R
	Add  = isa.Op{Name: "add", F7: 0b0000000, F3: 0b000, Functs: 0b0000000_000, Opcode: isa.Opcodes.R}
	Sub  = isa.Op{Name: "sub", F7: 0b0100000, F3: 0b000, Functs: 0b0100000_000, Opcode: isa.Opcodes.R}
	Sll  = isa.Op{Name: "sll", F7: 0b0000000, F3: 0b001, Functs: 0b0000000_001, Opcode: isa.Opcodes.R}
	Slt  = isa.Op{Name: "slt", F7: 0b0000000, F3: 0b010, Functs: 0b0000000_010, Opcode: isa.Opcodes.R}
	Sltu = isa.Op{Name: "sltu", F7: 0b0000000, F3: 011, Functs: 0b0000000_011, Opcode: isa.Opcodes.R}
	Xor  = isa.Op{Name: "xor", F7: 0b0000000, F3: 0b100, Functs: 0b0000000_100, Opcode: isa.Opcodes.R}
	Srl  = isa.Op{Name: "srl", F7: 0b0000000, F3: 0b101, Functs: 0b0000000_101, Opcode: isa.Opcodes.R}
	Sra  = isa.Op{Name: "sra", F7: 0b0100000, F3: 0b101, Functs: 0b0100000_101, Opcode: isa.Opcodes.R}
	Or   = isa.Op{Name: "or", F7: 0b0000000, F3: 0b110, Functs: 0b0000000_110, Opcode: isa.Opcodes.R}
	And  = isa.Op{Name: "and", F7: 0b0000000, F3: 0b111, Functs: 0b0000000_111, Opcode: isa.Opcodes.R}

	// RISC-32 Register-Immediate (RI)
	Addi  = isa.Op{Name: "addi", F7: 0b0000000, F3: 0b000, Functs: 0b0000000_000, Opcode: isa.Opcodes.RI}
	Slli  = isa.Op{Name: "slli", F7: 0b0000000, F3: 0b001, Functs: 0b0000000_001, Opcode: isa.Opcodes.RI}
	Slti  = isa.Op{Name: "slti", F7: 0b0000000, F3: 0b010, Functs: 0b0000000_010, Opcode: isa.Opcodes.RI}
	Sltiu = isa.Op{Name: "sltiu", F7: 0b0000000, F3: 0b011, Functs: 0b0000000_011, Opcode: isa.Opcodes.RI}
	Xori  = isa.Op{Name: "xori", F7: 0b0000000, F3: 0b100, Functs: 0b0000000_100, Opcode: isa.Opcodes.RI}
	Srli  = isa.Op{Name: "srli", F7: 0b0000000, F3: 0b101, Functs: 0b0000000_101, Opcode: isa.Opcodes.RI}
	Srai  = isa.Op{Name: "srai", F7: 0b0100000, F3: 0b101, Functs: 0b0100000_101, Opcode: isa.Opcodes.RI}
	Ori   = isa.Op{Name: "ori", F7: 0b0000000, F3: 0b110, Functs: 0b0000000_110, Opcode: isa.Opcodes.RI}
	Andi  = isa.Op{Name: "andi", F7: 0b0000000, F3: 0b111, Functs: 0b0000000_111, Opcode: isa.Opcodes.RI}

	// RISC-32 M
	Mul    = isa.Op{Name: "mul", F7: 0b0000001, F3: 0b000, Functs: 0b0000001_000, Opcode: isa.Opcodes.R}
	Mulh   = isa.Op{Name: "mulh", F7: 0b0000001, F3: 0b001, Functs: 0b0000001_001, Opcode: isa.Opcodes.R}
	Mulhsu = isa.Op{Name: "mulhsu", F7: 0b0000001, F3: 0b010, Functs: 0b0000001_010, Opcode: isa.Opcodes.R}
	Mulhu  = isa.Op{Name: "mulhu", F7: 0b0000001, F3: 0b011, Functs: 0b0000001_011, Opcode: isa.Opcodes.R}
	Div    = isa.Op{Name: "div", F7: 0b0000001, F3: 0b100, Functs: 0b0000001_100, Opcode: isa.Opcodes.R}
	Divu   = isa.Op{Name: "divu", F7: 0b0000001, F3: 0b101, Functs: 0b0000001_101, Opcode: isa.Opcodes.R}
	Rem    = isa.Op{Name: "rem", F7: 0b0000001, F3: 0b110, Functs: 0b0000001_110, Opcode: isa.Opcodes.R}
	Remu   = isa.Op{Name: "remu", F7: 0b0000001, F3: 0b111, Functs: 0b0000001_111, Opcode: isa.Opcodes.R}
)
