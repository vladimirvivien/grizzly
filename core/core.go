package core

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/ctrlunit"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

var (
	In = struct {
		Insts device.PinLabel
	}{
		Insts: "core.insts.in",
	}
)

type Core struct {
	*device.Base
	reg  device.Type
	alu  device.Type
	ctrl device.Type
}

func New() device.Type {
	return newCore()
}

func newCore() *Core {
	return &Core{
		Base: device.NewBase(),
		reg:  reg.New(),
		alu:  alu.New(),
		ctrl: ctrlunit.New(),
	}
}

// Run starts the core and its components
func (c *Core) Run() error {
	if err := c.wireComponents(); err != nil {
		return err
	}
	return c.startComponents()
}

func (c *Core) wireComponents() error {
	if c.GetPin(In.Insts) == nil {
		return fmt.Errorf("instructions datapath not set")
	}

	// wire controller
	c.ctrl.SetPin(ctrlunit.In.Insts, c.GetPin(In.Insts))

	// wire register
	c.reg.SetPin(reg.In.RDAddr, c.ctrl.GetPin(ctrlunit.Out.RD))
	c.reg.SetPin(reg.In.RS1Addr, c.ctrl.GetPin(ctrlunit.Out.RS1))
	c.reg.SetPin(reg.In.RS2Addr, c.ctrl.GetPin(ctrlunit.Out.RS2))
	c.reg.SetPin(reg.In.Werf, c.ctrl.GetPin(ctrlunit.Out.Werf))

	// wire alu
	c.alu.SetPin(alu.In.Operation, c.ctrl.GetPin(ctrlunit.Out.Functs))
	c.alu.SetPin(alu.In.Operand1, c.reg.GetPin(reg.Out.RS1Data))
	c.alu.SetPin(alu.In.Operand2, c.reg.GetPin(reg.Out.RS2Data))
	c.alu.SetPin(alu.In.Operand2, device.Select2(c.reg.GetPin(reg.Out.RS2Data), c.ctrl.GetPin(ctrlunit.Out.Imm)))
	c.reg.SetPin(reg.In.Data, c.alu.GetPin(alu.Out.Result))

	return nil
}

// startComponents loop through each component and invoke Run.
func (c *Core) startComponents() error {
	for _, comp := range []device.Type{c.ctrl, c.reg, c.alu} {
		if err := comp.Run(); err != nil {
			return err
		}
	}

	return nil
}
