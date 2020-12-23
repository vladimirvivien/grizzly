package isa

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

type Opcode = uint32

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
		S:     0x23,
		SB:    0x63,
		U:     0x37,
		UJ:    0x6F,
	}
)

// Inst represents a RISCV-32 instruction
type Inst = datapath.Word

func GetInstOpcode(i Inst) Opcode {
	return  i & 0x7F
}

// BaseFields common fields for all inst
type BaseFields struct {
	Opcode uint32 // opcode format
}

// IntegerBaseFields base for integer ops
type IntegerBaseFields struct {
	*BaseFields
	Rd     uint32 // destination register
	Funct3 uint32 // ISA 3-bit funct field
	Rs1    uint32 // register 1
}

func NewIntegerBaseFields() *IntegerBaseFields {
	return &IntegerBaseFields{BaseFields: &BaseFields{}}
}



// instruction operation
type Op struct {
	Name   string
	Functs uint32
	Opcode uint32
}


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
