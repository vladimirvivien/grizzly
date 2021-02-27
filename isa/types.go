package isa

type Opcode = uint8

var (
	Opcodes = struct {
		R,
		RI,
		L,
		Ecall,
		S,
		SB,
		U,
		UJ Opcode
	}{
		R:     0b0110011,
		RI:    0b0010011,
		L:     0b0000011,
		Ecall: 0x73,
		S:     0b0100011,
		SB:    0x63,
		U:     0x37,
		UJ:    0x6F,
	}
)

type Op struct {
	Name string
	Functs uint16
	Opcode uint8
}

type Fields struct {
	Opcode uint8
	Rd     uint8
	Funct3 uint8
	Rs1    uint8
	Rs2    uint8
	Funct7 uint8
	Shift  uint8
	Imm    uint32
}

// DecodeFuncts extracts both ISA funct fields funct7 and funct3 assuming
// value functs contain these values concatenated in the lower 10 bits as:
//
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
//
// If funct7 is not encoded, it's returned as zero.
func DecodeFuncts(functs uint8) (funct7, funct3 uint8) {
	funct3 = functs & 0x7
	funct7 = (functs >> 3) & 0x7F
	return
}
