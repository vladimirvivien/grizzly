package isa

var (
	Opcodes = struct {
		R,
		RI,
		L,
		Ecall,
		S,
		SB,
		U,
		UJ uint32
	}{
		R:     0b0110011,
		RI:    0b0010011,
		L:     0x03,
		Ecall: 0x73,
		S:     0x23,
		SB:    0x63,
		U:     0x37,
		UJ:    0x6F,
	}
)

// Inst represents a RISCV-32 instruction
type Inst = uint32

// BaseFields common fields for all inst
type BaseFields struct {
	Opcode uint32 // opcode format
}

// ComputeBaseFields base for integer ops
type ComputeBaseFields struct {
	*BaseFields
	Rd     uint32 // destination register
	Funct3 uint32 // ISA 3-bit funct field
	Rs1    uint32 // register 1
}

// RFields for integer register operations
type RFields struct {
	*ComputeBaseFields
	Rs2    uint32 // register 2
	Funct7 uint32 // ISA 7-bit funct field
}

// NewRFields returns &RFields{}
func NewRFields() *RFields {
	return &RFields{ComputeBaseFields: &ComputeBaseFields{BaseFields: &BaseFields{}}}
}

// Functs returns isa.Functs(f7,f3)
func (f *RFields) Functs() (result uint32) {
	return Functs(f.Funct7, f.Funct3)
}

// RIFields for Register-Immediate ops
type RIFields struct {
	*ComputeBaseFields
	Imm uint32
	// Shift when Funct3 = 001, 101
	Shift  uint32
	Funct7 uint32
}

// NewRIFields returns &RIFields{}
func NewRIFields() *RIFields {
	return &RIFields{ComputeBaseFields: &ComputeBaseFields{BaseFields: &BaseFields{}}}
}

// Functs returns isa.Functs(f7,f3)
func (f *RIFields) Functs() (result uint32) {
	return Functs(f.Funct7, f.Funct3)
}

// Functs encodes R-format functs fields funct7 and funct3
// as a single value into the lower 10 bits of result:
//
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
func Functs(f7, f3 uint32) (result uint32) {
	result = (result | f7) << 3
	result = result | f3
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

	// RISC-32 Register-Immediate (RI)
	Addi  = rop{Name: "addi", Functs: uint32(0b0000000_000), Opcode: Opcodes.RI}
	Slli  = rop{Name: "slli", Functs: uint32(0b0000000_001), Opcode: Opcodes.RI}
	Slti  = rop{Name: "slti", Functs: uint32(0b0000000_010), Opcode: Opcodes.RI}
	Sltiu = rop{Name: "sltiu", Functs: uint32(0b0000000_011), Opcode: Opcodes.RI}
	Xori  = rop{Name: "xori", Functs: uint32(0b0000000_100), Opcode: Opcodes.RI}
	Srli  = rop{Name: "srli", Functs: uint32(0b0000000_101), Opcode: Opcodes.RI}
	Srai  = rop{Name: "srai", Functs: uint32(0b0100000_101), Opcode: Opcodes.RI}
	Ori   = rop{Name: "ori", Functs: uint32(0b0000000_110), Opcode: Opcodes.RI}
	Andi  = rop{Name: "andi", Functs: uint32(0b0000000_111), Opcode: Opcodes.RI}

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
