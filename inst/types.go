package inst

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

type Type uint32

type Fields struct {
	Opcode uint32 // opcode format
	Rd     uint32 // destination register
	Funct3 uint32 // branch func
	Rs1    uint32 // register 1
	Rs2    uint32 // register 2
	Funct7 uint32 // alu func
	Imm    uint32 // immediate
}

// R-type operation
type rop struct {
	Name string
	// Funct encodes both funct7 + funct3 in its lower bits
	//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
	Funct  uint32
	Opcode uint32
}

// DecodeFuncts extracts both funct7 and funct3 from funct
// which is encoded as:
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
func (r rop) DecodeFuncts() (funct7, funct3 uint32) {
	funct3 = r.Funct & 0x7
	funct7 = (r.Funct >> 3) & 0x7F
	return
}

// EncRFuncts encodes R-type functions Funct7+Funct3 into the
// lower bits of a single value:
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
func EncRFuncts(f7, f3 uint32) uint32 {
	var result uint32
	result = (result | f7) << 3
	result = result | f3
	return result
}

// R-type operations
var (
	Funct7Lo uint32 = 0b0000000
	Funct7Hi uint32 = 0b0100000

	// Meta for R-Types
	Add  = rop{Name: "add", Funct: uint32(0) | 0b0000000000, Opcode: Opcodes.R}
	Sub  = rop{Name: "sub", Funct: uint32(0) | 0b0100000000, Opcode: Opcodes.R}
	Sll  = rop{Name: "sll", Funct: uint32(0) | 0b0000000001, Opcode: Opcodes.R}
	Slt  = rop{Name: "slt", Funct: uint32(0) | 0b0000000010, Opcode: Opcodes.R}
	Sltu = rop{Name: "sltu", Funct: uint32(0) | 0b0000000011, Opcode: Opcodes.R}
	Xor  = rop{Name: "xor", Funct: uint32(0) | 0b0000000100, Opcode: Opcodes.R}
	Srl  = rop{Name: "slr", Funct: uint32(0) | 0b0000000101, Opcode: Opcodes.R}
	Sra  = rop{Name: "sra", Funct: uint32(0) | 0b0100000101, Opcode: Opcodes.R}
	Or   = rop{Name: "or", Funct: uint32(0) | 0b0000000110, Opcode: Opcodes.R}
	And  = rop{Name: "and", Funct: uint32(0) | 0b0000000111, Opcode: Opcodes.R}
)
