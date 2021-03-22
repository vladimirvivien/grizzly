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

type Component interface {
	Run() error
}

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
		clock: clock.New(20*time.Microsecond),
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
	// wire decoder
	c.dec.Input(c.in)

	// wire register
	c.reg.OpInput(c.dec.Output())

	// wire alu
	c.alu.ParamsInput(c.reg.AluParamsOutput())

	// router
	c.rout.AluResultInput(c.alu.ResultOutput())
	c.reg.DataInput(c.rout.RegisterDataOutput())

}

func (c *Core) startComponents() error {
	comps := []Component{c.dec, c.reg, c.alu, c.rout}
	for _, comp := range comps {
		if err := comp.Run(); err != nil {
			return err
		}
	}
	return nil
}