package core

import (
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/decoder"
	"github.com/vladimirvivien/grizzly/mem"
	"github.com/vladimirvivien/grizzly/reg"
)

type Core struct {
	clock *clock.Clock
	in datapath.Bytestream
	dec *decoder.Decoder
	reg *reg.RegisterFile
	alu *alu.ALU
	mem *mem.Memory
}

func New() *Core {
	return &Core{
		clock: clock.New(2*time.Microsecond),
		dec: decoder.New(),
		reg: reg.New(),
		alu: alu.New(),
		mem: mem.New(1024*1000),
	}
}

func (c *Core) Input(in datapath.Bytestream) {
	c.in = in
}

func (c *Core) Run() error {
	c.wireDatapath()
	return c.startComponents()
}

func (c *Core) wireDatapath() {
	// decoder <- mem: instruction
	c.dec.Connect(decoder.Labels.Instruction, c.in)
	// reg <- decoder: op fields
	c.reg.Connect(reg.Labels.InFields, c.dec.GetPin(decoder.Labels.OutFields))
	// alu <- reg: Operation
	c.alu.Connect(alu.Labels.InOperations, c.reg.GetPin(reg.Labels.OutAluOps))
	// register <- alu: register data
	c.reg.Connect(reg.Labels.InAluData, c.alu.GetPin(alu.Labels.OutRegData))
	// mem <- alu: mem op
	c.mem.Connect(mem.Labels.InOperation, c.alu.GetPin(alu.Labels.OutMemOp))
	// register <- mem: register data
	c.reg.Connect(reg.Labels.InMemData, c.mem.GetPin(mem.Labels.OutRegData))
}

func (c *Core) startComponents() error {
	components := []datapath.Component{c.dec, c.reg, c.alu, c.mem}
	for _, comp := range components {
		if err := comp.Run(); err != nil {
			return err
		}
	}
	return nil
}