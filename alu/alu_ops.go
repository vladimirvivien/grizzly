package alu

import (
	"github.com/vladimirvivien/grizzly/isa/integer"
)

var (
	Ops = struct {
		Add  uint8
		And  uint8
		Sub  uint8
		Or   uint8
		Sll  uint8
		Slt  uint8
		Sltu uint8
		Sra  uint8
		Srl  uint8
		Xor  uint8

		// mul
		Mul    uint8
		Mulh   uint8
		Mulhsu uint8
		Mulhu  uint8

		// div
		Div  uint8
		Divu uint8
		Rem  uint8
		Remu uint8

		// Branch
		Branch0  uint8
		Branch1 uint8
	}{
		Add:  0b00000000,
		And:  0b00000001,
		Sub:  0b00000010,
		Or:   0b00000011,
		Sll:  0b00000100,
		Slt:  0b00000101,
		Sltu: 0b00000110,
		Sra:  0b00000111,
		Srl:  0b00001000,
		Xor:  0b00001001,

		// Mul
		Mul:    0b00001010,
		Mulh:   0b00001011,
		Mulhsu: 0b00001100,
		Mulhu:  0b00001101,

		// Div
		Div:  0b00001110,
		Divu: 0b00001111,
		Rem:  0b00010000,
		Remu: 0b00010001,

		// comp
		Branch0:  0b01000000,
		Branch1:  0b01000001,
	}
)

// EncodeAluOp returns a canonical operation code
// for the ALU based on the functs values provided
func EncodeAluOp(f7, f3 uint8) uint8 {
	switch {
	case integer.Add.F7 == f7 && integer.Add.F3 == f3,
		integer.Addi.F7 == f7 && integer.Addi.F3 == f3:
		return Ops.Add

	case integer.Sub.F7 == f7 && integer.Sub.F3 == f3:
		return Ops.Sub

	case integer.Sll.F7 == f7 && integer.Sll.F3 == f3,
		integer.Slli.F7 == f7 && integer.Slli.F3 == f3:
		return Ops.Sll

	case integer.Slt.F7 == f7 && integer.Slt.F3 == f3,
		integer.Slti.F7 == f7 && integer.Slti.F3 == f3:
		return Ops.Slt

	case integer.Sltu.F7 == f7 && integer.Sltu.F3 == f3,
		integer.Sltiu.F7 == f7 && integer.Sltiu.F3 == f3:
		return Ops.Sltu

	case integer.Xor.F7 == f7 && integer.Xor.F3 == f3,
		integer.Xori.F7 == f7 && integer.Xori.F3 == f3:
		return Ops.Xor

	case integer.Srl.F7 == f7 && integer.Srl.F3 == f3,
		integer.Srli.F7 == f7 && integer.Srli.F3 == f3:
		return Ops.Srl

	case integer.Sra.F7 == f7 && integer.Sra.F3 == f3,
		integer.Srai.F7 == f7 && integer.Srai.F3 == f3:
		return Ops.Sra

	case integer.Or.F7 == f7 && integer.Or.F3 == f3,
		integer.Ori.F7 == f7 && integer.Ori.F3 == f3:
		return Ops.Or

	case integer.And.F7 == f7 && integer.And.F3 == f3,
		integer.Andi.F7 == f7 && integer.Andi.F3 == f3:
		return Ops.And

	case integer.Mul.F7 == f7 && integer.Mul.F3 == f3:
		return Ops.Mul

	case integer.Mulh.F7 == f7 && integer.Mulh.F3 == f3:
		return Ops.Mulh

	case integer.Mulhsu.F7 == f7 && integer.Mulhsu.F3 == f3:
		return Ops.Mulhsu

	case integer.Mulhu.F7 == f7 && integer.Mulhu.F3 == f3:
		return Ops.Mulhu

	case integer.Div.F7 == f7 && integer.Div.F3 == f3:
		return Ops.Div

	case integer.Divu.F7 == f7 && integer.Divu.F3 == f3:
		return Ops.Divu

	case integer.Rem.F7 == f7 && integer.Rem.F3 == f3:
		return Ops.Rem

	case integer.Remu.F7 == f7 && integer.Remu.F3 == f3:
		return Ops.Remu

	default:
		panic("unknown instruction functs")
	}
}