package isa

var (
	Opcodes = struct {
		R,
		L,
		RI,
		Ecall,
		S,
		SB,
		U,
		UJ uint32
	}{
		R:     0b0110011,
		L:     0x03,
		RI:    0x13,
		Ecall: 0x73,
		S:     0x23,
		SB:    0x63,
		U:     0x37,
		UJ:    0x6F,
	}
)

// Inst represents a RISCV-32 instruction
type Inst = uint32

// RFields represents instruction fields for R-format
type RFields struct {
	Opcode uint32 // opcode format
	Rd     uint32 // destination register
	Funct3 uint32 // ISA 3-bit funct field
	Rs1    uint32 // register 1
	Rs2    uint32 // register 2
	Funct7 uint32 // ISA 7-bit funct field
}

// Functs encodes R-format fields funct7 and funct3 by concatenating
// their values into the lower 10 bits of result:
//
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
func (f *RFields) Functs() (result uint32) {
	result = (result | f.Funct7) << 3
	result = result | f.Funct3
	return result
}

// R-type operation
type rop struct {
	Name   string
	Functs uint32
	Opcode uint32
}

var (
	// RISC-32 R
	Add  = rop{Name: "add", Functs: uint32(0b0000000_000), Opcode: Opcodes.R}
	Sub  = rop{Name: "sub", Functs: uint32(0b0100000_000), Opcode: Opcodes.R}
	Sll  = rop{Name: "sll", Functs: uint32(0b0000000_001), Opcode: Opcodes.R}
	Slt  = rop{Name: "slt", Functs: uint32(0b0000000_010), Opcode: Opcodes.R}
	Sltu = rop{Name: "sltu", Functs: uint32(0b0000000_011), Opcode: Opcodes.R}
	Xor  = rop{Name: "xor", Functs: uint32(0b0000000_100), Opcode: Opcodes.R}
	Srl  = rop{Name: "srl", Functs: uint32(0b0000000_101), Opcode: Opcodes.R}
	Sra  = rop{Name: "sra", Functs: uint32(0b0100000_101), Opcode: Opcodes.R}
	Or   = rop{Name: "or", Functs: uint32(0b0000000_110), Opcode: Opcodes.R}
	And  = rop{Name: "and", Functs: uint32(0b0000000_111), Opcode: Opcodes.R}

	// RISC-32 M
	Mul    = rop{Name: "mul", Functs: uint32(0b0000001_000), Opcode: Opcodes.R}
	Mulh   = rop{Name: "mulh", Functs: uint32(0b0000001_001), Opcode: Opcodes.R}
	Mulhsu = rop{Name: "mulhsu", Functs: uint32(0b0000001_010), Opcode: Opcodes.R}
	Mulhu  = rop{Name: "mulhu", Functs: uint32(0b0000001_011), Opcode: Opcodes.R}
	Div    = rop{Name: "div", Functs: uint32(0b0000001_100), Opcode: Opcodes.R}
	Divu   = rop{Name: "divu", Functs: uint32(0b0000001_101), Opcode: Opcodes.R}
	Rem    = rop{Name: "rem", Functs: uint32(0b0000001_110), Opcode: Opcodes.R}
	Remu   = rop{Name: "remu", Functs: uint32(0b0000001_111), Opcode: Opcodes.R}
)

// DecFuncts extracts both ISA funct fields funct7 and funct3 assuming
// value functs contain these values concatenated in the lower 10 bits as:
//
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
//
// If funct7 is not encoded, it's returned as zero.
func DecFuncts(functs uint32) (funct7, funct3 uint32) {
	funct3 = functs & 0x7
	funct7 = (functs >> 3) & 0x7F
	return
}
