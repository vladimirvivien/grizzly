package brancher

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/branch"
)

var (
	Labels = struct {
		InBranchOp datapath.Pin
		OutOperation  datapath.Pin
	}{
		InBranchOp: datapath.Pin("branch.in.op"),
	}
)

type Brancher struct {
	*datapath.BaseComponent
	outOperation chan []byte
}

func New() *Brancher {
	br := &Brancher{
		BaseComponent: datapath.NewBase(),
		outOperation: make(chan []byte),
	}
	br.Connect(Labels.OutOperation, br.outOperation)
	return br
}

func (b *Brancher) Run() error {
	input := b.GetPin(Labels.InBranchOp)
	if input == nil {
		return fmt.Errorf("brancher: missing input: %s", Labels.InBranchOp)
	}
	go func () {
		defer close(b.outOperation)
		for stream := range input{
			brOp := datapath.DecodeBranchOp(stream)
			aluOp := datapath.Operation{
				Opcode:      brOp.Opcode,
			}

			switch brOp.Funct3{
			case branch.Beq.F3:
				if datapath.SXWord(brOp.RS1D) == datapath.SXWord(brOp.RS1D) {
					aluOp.AluOp = alu.Ops.Branch1
					aluOp.AluOperand1 = brOp.PC
					aluOp.AluOperand2 = brOp.Imm
				}
			case branch.Bne.F3:
			case branch.Blt.F3:
			case branch.Bge.F3:
			case branch.Bltu.F3:
			case branch.Bgeu.F3:
			}

			b.outOperation <- datapath.EncodeOp(aluOp)
		}

	}()
return nil
}
