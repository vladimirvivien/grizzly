package datapath

import (
	"encoding/binary"
)

// TODO
// Investigate using []byte to represent instruction pins between components. This
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
const (
	RegSize    = 32
	XWordLen   = 32
	XWordBytes = XWordLen / 8
)

// XWord represents the appropriate
// Word size for given ISA spec (16,32,64,128 bits)
type XWord = uint32
type SXWord = int32

// ProgramCounter represents the resolved PC value that can be
// used as address the next program address to load.
type ProgramCounter = XWord

func DecodePC(stream []byte) ProgramCounter {
	return binary.LittleEndian.Uint32(stream)
}

func EncodePC(pc ProgramCounter) []byte {
	buf := make([]byte, XWordBytes, XWordBytes)
	binary.LittleEndian.PutUint32(buf, pc)
	return buf
}

// PcOp represents an operation that can set the value of the
// next selected program counter value if Jump > 0. The bytestream
// for PcOp is shown below:
//
// 0       1       2       3       4       5
// 01234567012345670123456701234567012345670
// +-------+-------+-------+-------+-------+
// |  Jump |              Ins              |
// +-------+-------+-------+-------+-------+
//
type PcOp struct {
	Jump uint8
	PC XWord
}

func DecodePcOp(stream []byte) PcOp {
	return PcOp{
		Jump:   stream[0],
		PC: binary.LittleEndian.Uint32(stream[1:]),
	}
}

func EncodePcOp(op PcOp) []byte {
	buf := make([]byte, 5, 5)
	buf[0] = op.Jump
	binary.LittleEndian.PutUint32(buf[1:], op.PC)
	return buf
}


// Instruction represents the instruction value at the
// specified PC from the instruction memory. The bytestream
// for Instruction is shown below:
//
//
// 0       1       2       3       4       5       6       7
// 0123456701234567012345670123456701234567012345670123456701234567
// +-------+-------+-------+-------+-------+-------+------+-------+
// |               PC              |              Ins             |
// +-------+-------+-------+-------+-------+-------+------+-------+
//
//
type Instruction struct {
	PC   XWord
	Inst XWord
}

func DecodeInstruction(stream []byte) Instruction {
	return Instruction{
		PC:   binary.LittleEndian.Uint32(stream[0:]),
		Inst: binary.LittleEndian.Uint32(stream[4:]),
	}
}

func EncodeInstruction(ins Instruction) []byte {
	buf := make([]byte, 8, 8)
	binary.LittleEndian.PutUint32(buf[0:], ins.PC)
	binary.LittleEndian.PutUint32(buf[4:], ins.Inst)
	return buf
}

// OpFields represent the operation and data fields decoded from the instruction.
// The bytestream from the decoder will have the following layout:
//
// 0       1       2       3       4       5       6       7
// 01234567012345670123456701234567012345670123456701234567012345670
// +-------+-------+-------+-------+-------+-------+-------+-------+
// |               PC              |OpCode |   Rd  |Funct3 |   Rs1 |
// +-------+-------+-------+-------+-------+-------+-------+-------+
// |  Rs2  |Funct7 |Shift  |              Imm              |
// +-----------------------+-------------------------------+
//
type OpFields struct {
	PC     XWord
	Opcode uint8
	Rd     uint8
	Funct3 uint8
	Rs1    uint8
	Rs2    uint8
	Funct7 uint8
	Shift  uint8
	Imm    uint32
}

func DecodeOpFields(stream []byte) OpFields {
	return OpFields{
		PC:     binary.LittleEndian.Uint32(stream[0:]),
		Opcode: stream[4],
		Rd:     stream[5],
		Funct3: stream[6],
		Rs1:    stream[7],
		Rs2:    stream[8],
		Funct7: stream[9],
		Shift:  stream[10],
		Imm:    binary.LittleEndian.Uint32(stream[11:]),
	}
}

func EncodeOpFields(f OpFields) []byte {
	buf := make([]byte, 15, 15)
	binary.LittleEndian.PutUint32(buf[0:], f.PC)
	buf[4] = f.Opcode
	buf[5] = f.Rd
	buf[6] = f.Funct3
	buf[7] = f.Rs1
	buf[8] = f.Rs2
	buf[9] = f.Funct7
	buf[10] = f.Shift
	binary.LittleEndian.PutUint32(buf[11:], f.Imm)
	return buf
}

// Operation represents the operation to be carried out in the ALU.
// The bytestrem to the ALU is encoded with the following layout:
//
// 0       1       2       3       4       5       6       7
// 0123456701234567012345670123456701234567012345670123456701234567
// +-------+-------+-------+-------+-------+-------+------+-------+
// |OpCode |   Rd  | AluOp |           AluOperand1        |
// +-------+-------+-------+-------+-------+-------+------+-------+
//        AluOperand2      | MemOp |              MemData         |
// +-------+-----------------------+------------------------------+
//
//
type Operation struct {
	// inst fields
	Opcode,

	// register store
	Rd uint8

	// ALU operation
	AluOp uint8
	AluOperand1,
	AluOperand2 XWord

	// Mem operation
	MemOp   uint8
	MemData XWord
}

func DecodeOp(s []byte) Operation {
	return Operation{
		Opcode:      s[0],
		Rd:          s[1],
		AluOp:       s[2],
		AluOperand1: binary.LittleEndian.Uint32(s[3:]),
		AluOperand2: binary.LittleEndian.Uint32(s[7:]),
		MemOp:       s[11],
		MemData:     binary.LittleEndian.Uint32(s[12:]),
	}
}

func EncodeOp(a Operation) []byte {
	buf := make([]byte, 16, 16)
	buf[0] = a.Opcode
	buf[1] = a.Rd
	buf[2] = a.AluOp
	binary.LittleEndian.PutUint32(buf[3:], a.AluOperand1)
	binary.LittleEndian.PutUint32(buf[7:], a.AluOperand2)
	buf[11] = a.MemOp
	binary.LittleEndian.PutUint32(buf[12:], a.MemData)
	return buf
}

// RegisterData represents data to be stored in the register at the end of an operation.
// The bytestrem for register data is encoded with the following layout:
//
// 0       1       2       3       4
// 0123456701234567012345670123456701234567
// +-------+-------+-------+-------+-------+
// |Rd     |              Value            |
// +-------+-------+-------+-------+-------+
//
type RegisterData struct {
	Rd    uint8
	Value XWord
}

func DecodeRegStore(s []byte) RegisterData {
	return RegisterData{
		Rd:    s[0],
		Value: binary.LittleEndian.Uint32(s[1:]),
	}
}

func EncodeRegStore(r RegisterData) []byte {
	buf := make([]byte, 5, 5)
	buf[0] = r.Rd
	binary.LittleEndian.PutUint32(buf[1:], r.Value)
	return buf
}

// MemOp represents memory operation directives for the memory component.
// The bytestrem for the operations is encoded with the following layout:
//
// 0       1       2       3       4       5       6       7
// 0123456701234567012345670123456701234567012345670123456701234567
// +-------+-------+-------+-------+-------+-------+------+-------+
// |OpCode |   Rd  |   Op  |              Addr            |
// +-------+-------+-------+-------+-------+-------+------+-------+
//            Data         |
// +-------+-------+-------+
type MemOp struct {
	Opcode,

	// reg store
	Rd,

	// memory op params
	Op uint8
	Addr XWord
	Data XWord
}

func DecodeMemOp(s []byte) MemOp {
	return MemOp{
		Opcode: s[0],
		Rd:     s[1],
		Op:     s[2],
		Addr:   binary.LittleEndian.Uint32(s[3:]),
		Data:   binary.LittleEndian.Uint32(s[7:]),
	}
}

func EncodeMemOp(o MemOp) []byte {
	buf := make([]byte, 12, 12)
	buf[0] = o.Opcode
	buf[1] = o.Rd
	buf[2] = o.Op
	binary.LittleEndian.PutUint32(buf[3:], o.Addr)
	binary.LittleEndian.PutUint32(buf[7:], o.Data)
	return buf
}
