package ctrlunit

import (
	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/isa"
)

func encodeAluOp(functs uint32) uint32 {
	switch functs {
	case isa.Add.Functs, isa.Addi.Functs:
		return alu.Ops.Add
	case isa.Sub.Functs:
		return alu.Ops.Sub
	case isa.Sll.Functs, isa.Slli.Functs:
		return alu.Ops.Sll
	case isa.Slt.Functs, isa.Slti.Functs:
		return alu.Ops.Slt
	case isa.Sltu.Functs, isa.Sltiu.Functs:
		return alu.Ops.Sltu
	case isa.Xor.Functs, isa.Xori.Functs:
		return alu.Ops.Xor
	case isa.Srl.Functs, isa.Srli.Functs:
		return alu.Ops.Srl
	case isa.Sra.Functs, isa.Srai.Functs:
		return alu.Ops.Sra
	case isa.Or.Functs, isa.Ori.Functs:
		return alu.Ops.Or
	case isa.And.Functs, isa.Andi.Functs:
		return alu.Ops.And
	default:
		return 0
	}
}