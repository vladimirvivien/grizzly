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
		UJ,
		J,
		JI Opcode
	}{
		R:     0b0110011,
		RI:    0b0010011,
		L:     0b0000011,
		Ecall: 0x73,
		S:     0b0100011,
		SB:    0x63,
		U:     0x37,
		UJ:    0x6F,
		J:     0b1101111,
		JI:    0b1100111,
	}
)

type Op struct {
	Name   string
	Functs uint16
	F3, F7 uint8
	Opcode uint8
}
