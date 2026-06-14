//go:build rv64 || rv64i

package datapath

import (
	"encoding/binary"
)

const (
	RegSize    = 64
	XWordLen   = 64
	XWordBytes = XWordLen / 8
)

// XWord represents the appropriate
// Word size for given ISA spec (16,32,64,128 bits)
type XWord = uint64
type SXWord = int64

func DecodeXWord(stream []byte) XWord {
	return binary.LittleEndian.Uint64(stream)
}

func EncodeXWord(w XWord) []byte {
	buf := make([]byte, XWordBytes, XWordBytes)
	binary.LittleEndian.PutUint64(buf, w)
	return buf
}

// ProgramCounter represents the resolved PC value that can be
// used as address the next program address to load.
type ProgramCounter = XWord

func DecodePC(stream []byte) ProgramCounter {
	return binary.LittleEndian.Uint64(stream)
}

func EncodePC(pc ProgramCounter) []byte {
	buf := make([]byte, XWordBytes, XWordBytes)
	binary.LittleEndian.PutUint64(buf, pc)
	return buf
}

// PcOp represents an operation that can set the value of the
// next selected program counter value if Jump > 0.
type PcOp struct {
	Jump uint8
	PC   XWord
}

func DecodePcOp(stream []byte) PcOp {
	return PcOp{
		Jump: stream[0],
		PC:   binary.LittleEndian.Uint64(stream[1:]),
	}
}

func EncodePcOp(op PcOp) []byte {
	buf := make([]byte, 1+XWordBytes, 1+XWordBytes)
	buf[0] = op.Jump
	binary.LittleEndian.PutUint64(buf[1:], op.PC)
	return buf
}

// Instruction represents the instruction value at the
// specified PC from the instruction memory.
type Instruction struct {
	PC   XWord
	Inst uint32
}

func DecodeInstruction(stream []byte) Instruction {
	return Instruction{
		PC:   binary.LittleEndian.Uint64(stream[0:]),
		Inst: binary.LittleEndian.Uint32(stream[8:]),
	}
}

func EncodeInstruction(ins Instruction) []byte {
	buf := make([]byte, 12, 12)
	binary.LittleEndian.PutUint64(buf[0:], ins.PC)
	binary.LittleEndian.PutUint32(buf[8:], ins.Inst)
	return buf
}

// OpFields represent the operation and data fields decoded from the instruction.
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
		PC:     binary.LittleEndian.Uint64(stream[0:]),
		Opcode: stream[8],
		Rd:     stream[9],
		Funct3: stream[10],
		Rs1:    stream[11],
		Rs2:    stream[12],
		Funct7: stream[13],
		Shift:  stream[14],
		Imm:    binary.LittleEndian.Uint32(stream[15:]),
	}
}

func EncodeOpFields(f OpFields) []byte {
	buf := make([]byte, 19, 19)
	binary.LittleEndian.PutUint64(buf[0:], f.PC)
	buf[8] = f.Opcode
	buf[9] = f.Rd
	buf[10] = f.Funct3
	buf[11] = f.Rs1
	buf[12] = f.Rs2
	buf[13] = f.Funct7
	buf[14] = f.Shift
	binary.LittleEndian.PutUint32(buf[15:], f.Imm)
	return buf
}

// Operation represents the operation to be carried out in the ALU.
type Operation struct {
	PC     XWord
	Opcode uint8
	Rd     uint8
	AluOp  uint8
	AluOperand1,
	AluOperand2 XWord
	MemOp   uint8
	MemData XWord
}

func DecodeOp(s []byte) Operation {
	return Operation{
		PC:          binary.LittleEndian.Uint64(s[0:]),
		Opcode:      s[8],
		Rd:          s[9],
		AluOp:       s[10],
		AluOperand1: binary.LittleEndian.Uint64(s[11:]),
		AluOperand2: binary.LittleEndian.Uint64(s[19:]),
		MemOp:       s[27],
		MemData:     binary.LittleEndian.Uint64(s[28:]),
	}
}

func EncodeOp(a Operation) []byte {
	buf := make([]byte, 36, 36)
	binary.LittleEndian.PutUint64(buf[0:], a.PC)
	buf[8] = a.Opcode
	buf[9] = a.Rd
	buf[10] = a.AluOp
	binary.LittleEndian.PutUint64(buf[11:], a.AluOperand1)
	binary.LittleEndian.PutUint64(buf[19:], a.AluOperand2)
	buf[27] = a.MemOp
	binary.LittleEndian.PutUint64(buf[28:], a.MemData)
	return buf
}

// BranchOp carries control/data for branch operation to be carried
// out by the Brancher.
type BranchOp struct {
	PC     XWord
	Opcode uint8
	Funct3 uint8
	RS1D   XWord
	RS2D   XWord
	Imm    uint32
}

func DecodeBranchOp(s []byte) BranchOp {
	return BranchOp{
		PC:     binary.LittleEndian.Uint64(s[0:]),
		Opcode: s[8],
		Funct3: s[9],
		RS1D:   binary.LittleEndian.Uint64(s[10:]),
		RS2D:   binary.LittleEndian.Uint64(s[18:]),
		Imm:    binary.LittleEndian.Uint32(s[26:]),
	}
}

func EncodeBranchOp(a BranchOp) []byte {
	buf := make([]byte, 30, 30)
	binary.LittleEndian.PutUint64(buf[0:], a.PC)
	buf[8] = a.Opcode
	buf[9] = a.Funct3
	binary.LittleEndian.PutUint64(buf[10:], a.RS1D)
	binary.LittleEndian.PutUint64(buf[18:], a.RS2D)
	binary.LittleEndian.PutUint32(buf[26:], a.Imm)
	return buf
}

// RegisterData represents data to be stored in the register at the end of an operation.
type RegisterData struct {
	Rd    uint8
	Value XWord
}

func DecodeRegData(s []byte) RegisterData {
	return RegisterData{
		Rd:    s[0],
		Value: binary.LittleEndian.Uint64(s[1:]),
	}
}

func EncodeRegData(r RegisterData) []byte {
	buf := make([]byte, 9, 9)
	buf[0] = r.Rd
	binary.LittleEndian.PutUint64(buf[1:], r.Value)
	return buf
}

// MemOp represents memory operation directives for the memory component.
type MemOp struct {
	Opcode uint8
	Rd     uint8
	Op     uint8
	Addr   XWord
	Data   XWord
}

func DecodeMemOp(s []byte) MemOp {
	return MemOp{
		Opcode: s[0],
		Rd:     s[1],
		Op:     s[2],
		Addr:   binary.LittleEndian.Uint64(s[3:]),
		Data:   binary.LittleEndian.Uint64(s[11:]),
	}
}

func EncodeMemOp(o MemOp) []byte {
	buf := make([]byte, 19, 19)
	buf[0] = o.Opcode
	buf[1] = o.Rd
	buf[2] = o.Op
	binary.LittleEndian.PutUint64(buf[3:], o.Addr)
	binary.LittleEndian.PutUint64(buf[11:], o.Data)
	return buf
}
