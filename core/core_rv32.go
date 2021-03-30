package core

import (
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/decoder"
	"github.com/vladimirvivien/grizzly/reg"
	"github.com/vladimirvivien/grizzly/router"
)

type Core struct {
	clock *clock.Clock
	in datapath.Bytestream
	dec *decoder.Decoder
	reg *reg.RegisterFile
	alu *alu.ALU
	rout *router.Router
}

func New() *Core {
	return &Core{
		clock: clock.New(2*time.Microsecond),
		dec: decoder.New(),
		reg: reg.New(),
		alu: alu.New(),
		rout: router.New(),
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
	c.dec.Connect(decoder.Labels.Instruction, c.in)
	c.reg.Connect(reg.Labels.InFields, c.dec.GetPin(decoder.Labels.OutFields))
	c.alu.Connect(alu.Labels.InParams, c.reg.GetPin(reg.Labels.OutAluParams))
	c.rout.Connect(router.Labels.InAluResult, c.alu.GetPin(alu.Labels.OutResult))
	c.reg.Connect(reg.Labels.InData, c.rout.GetPin(router.Labels.OutRegisterData))
}

func (c *Core) startComponents() error {
	comps := []datapath.Component{c.dec, c.reg, c.alu, c.rout}
	for _, comp := range comps {
		if err := comp.Run(); err != nil {
			return err
		}
	}
	return nil
}