package brancher

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/branch"
)

var (
	Labels = struct {
		InBranchOp   datapath.Pin
		OutOperation datapath.Pin
	}{
		InBranchOp:   datapath.Pin("branch.in.op"),
		OutOperation: datapath.Pin("branch.out.operation"),
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
	go func() {
		defer close(b.outOperation)
		for stream := range input {
			brOp := datapath.DecodeBranchOp(stream)
			
			// Reconstruct the actual signed offset from the decoded Imm
			// brOp.Imm contains the 12-bit un-shifted, unsigned immediate value.
			offset := brOp.Imm << 1
			if (offset & 0x1000) != 0 { // if bit 12 is set, sign-extend
				offset |= 0xffffe000
			}

			aluOp := datapath.Operation{
				Opcode:      brOp.Opcode,
				AluOp:       alu.Ops.Branch1,
				AluOperand1: brOp.PC,
				AluOperand2: 4, // default: branch not taken, PC + 4
			}

			taken := false
			switch brOp.Funct3 {
			case branch.Beq.F3:
				taken = brOp.RS1D == brOp.RS2D
			case branch.Bne.F3:
				taken = brOp.RS1D != brOp.RS2D
			case branch.Blt.F3:
				taken = datapath.SXWord(brOp.RS1D) < datapath.SXWord(brOp.RS2D)
			case branch.Bge.F3:
				taken = datapath.SXWord(brOp.RS1D) >= datapath.SXWord(brOp.RS2D)
			case branch.Bltu.F3:
				taken = brOp.RS1D < brOp.RS2D
			case branch.Bgeu.F3:
				taken = brOp.RS1D >= brOp.RS2D
			}

			if taken {
				aluOp.AluOperand2 = offset
			}

			b.outOperation <- datapath.EncodeOp(aluOp)
		}
	}()
	return nil
}
