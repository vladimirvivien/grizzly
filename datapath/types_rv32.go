package datapath

// TODO
// Investigate using []byte to represent instruction streams between components. This
// approach would allow Grizzly to support multi-width instructions (32, 64, 128, etc).
//
// How would this work?
//   - The Wires type would be defined as chan []byte to stream both control/program data.
//   - Instructions would be converted from a numeric value to a stream of bytes
//   - The stream of bytes would be sent to core components over Wires
//   - When components get []byte, the bytes are converted to numeric representation
//     using the encoding/binary package to narrow to a specific value based on configured
//     instruction width.
//
// See https://play.golang.org/p/py_Uv9zSXWv
//
// This change would allow Grizzly to support different implementations of the RISCV
// ISA including compressed instructions based on the size of XLEN.
const XLEN = 32
const Width32 = 32
const Width64 = 64

// XWord represents the appropriate
// Word size for given ISA spec (16,32,64,128 bits)
type XWord = uint32

// OpFields operation control and data
// constructed from instruction.
type OpFields struct {
	Opcode uint8
	Rd     uint8
	Funct3 uint8
	Rs1    uint8
	Rs2    uint8
	Funct7 uint8
	Shift  uint8
	Imm    uint32
}

type AluParam struct {
	Opcode,
	Rd,
	Funct3,
	Funct7 uint8

	Op1,
	Op2 XWord
}

type AluResult struct {
	F3 uint8
	Rd uint8
	Value XWord
}

type RegisterData struct {
	Rd uint8
	Value XWord
}
