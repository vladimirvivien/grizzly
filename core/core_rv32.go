package core

import (
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/decoder"
	"github.com/vladimirvivien/grizzly/mem/data"
	"github.com/vladimirvivien/grizzly/mem/instruction"
	"github.com/vladimirvivien/grizzly/pc"
	"github.com/vladimirvivien/grizzly/reg"
)

type Core struct {
	clock *clock.Clock
	pc    *pc.PC
	imem  *instruction.InstructionMemory

	dec   *decoder.Decoder
	reg   *reg.RegisterFile
	alu   *alu.ALU
	dmem  *data.DataMemory
}

func New() *Core {
	return &Core{
		clock: clock.New(2 * time.Microsecond),
		pc:    pc.New(),
		imem: instruction.New(1024),
		dec:   decoder.New(),
		reg:   reg.New(),
		alu:   alu.New(),
		dmem:  data.New(1024 * 1000),
	}
}

func (c *Core) Run() error {
	c.wireAll()
	return c.startComponents()
}

func (c *Core) wireAll() {
	// inst mem <- pc
	c.imem.Connect(instruction.Labels.InPC, c.pc.GetPin(pc.Labels.OutCounter))
	// decoder <- dmem: instruction
	c.dec.Connect(decoder.Labels.Instruction, c.imem.GetPin(instruction.Labels.OutInstruction))
	// reg <- decoder: op fields
	c.reg.Connect(reg.Labels.InFields, c.dec.GetPin(decoder.Labels.OutFields))
	// alu <- reg: Operation
	c.alu.Connect(alu.Labels.InOperations, c.reg.GetPin(reg.Labels.OutAluOps))
	// register <- alu: register data
	c.reg.Connect(reg.Labels.InAluData, c.alu.GetPin(alu.Labels.OutRegData))
	// dmem <- alu: dmem op
	c.dmem.Connect(data.Labels.InOperation, c.alu.GetPin(alu.Labels.OutMemOp))
	// register <- dmem: register data
	c.reg.Connect(reg.Labels.InMemData, c.dmem.GetPin(data.Labels.OutRegData))
	// pc <- alu: PC op
	c.pc.Connect(pc.Labels.InPcOp, c.alu.GetPin(alu.Labels.OutPcOp))
}

func (c *Core) startComponents() error {
	components := []datapath.Component{c.pc, c.imem, c.dec, c.reg, c.alu, c.dmem}
	for _, comp := range components {
		if err := comp.Run(); err != nil {
			return err
		}
	}
	return nil
}
