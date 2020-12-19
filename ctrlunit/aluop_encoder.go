package ctrlunit

import (
	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

func encodeAluOp(functs uint32) uint32 {
	switch functs {
	case integer.Add.Functs, integer.Addi.Functs:
		return alu.Ops.Add
	case integer.Sub.Functs:
		return alu.Ops.Sub
	case integer.Sll.Functs, integer.Slli.Functs:
		return alu.Ops.Sll
	case integer.Slt.Functs, integer.Slti.Functs:
		return alu.Ops.Slt
	case integer.Sltu.Functs, integer.Sltiu.Functs:
		return alu.Ops.Sltu
	case integer.Xor.Functs, integer.Xori.Functs:
		return alu.Ops.Xor
	case integer.Srl.Functs, integer.Srli.Functs:
		return alu.Ops.Srl
	case integer.Sra.Functs, integer.Srai.Functs:
		return alu.Ops.Sra
	case integer.Or.Functs, integer.Ori.Functs:
		return alu.Ops.Or
	case integer.And.Functs, integer.Andi.Functs:
		return alu.Ops.And
	default:
		return 0
	}
}