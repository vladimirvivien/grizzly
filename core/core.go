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
	pins         device.Pins
	instructions device.WiresIn
	reg          device.Type
	alu          device.Type
	ctrl         device.Type
}

func New() device.Type {
	return newCore()
}

func newCore() *Core {
	return &Core{
		reg:  reg.New(),
		alu:  alu.New(),
		ctrl: ctrlunit.New(),
		pins: make(device.Pins),
	}
}

func (c *Core) GetPins() device.Pins {
	return c.pins
}

func (c *Core) GetPin(label device.PinLabel) device.Pin {
	return c.pins[label]
}

func (c *Core) SetPin(label device.PinLabel, pin device.Pin) {
	c.pins[label] = pin
}

func (c *Core) Run() error {
	if err := c.wireComponents(); err != nil {
		return err
	}
	return c.startComponents()
}

func (c *Core) wireComponents() error {
	if c.pins[In.Insts] == nil {
		return fmt.Errorf("instructions datapath not set")
	}

	// wire controller
	c.ctrl.SetPin(ctrlunit.In.Insts, c.instructions)

	// wire register
	c.reg.SetPin(reg.In.RS1Addr, c.ctrl.GetPin(ctrlunit.Out.RS1))
	c.reg.SetPin(reg.In.RS2Addr, c.ctrl.GetPin(ctrlunit.Out.RS2))
	c.reg.SetPin(reg.In.RDAddr, c.ctrl.GetPin(ctrlunit.Out.RD))

	// wire alu
	c.alu.SetPin(alu.In.Functs, c.ctrl.GetPin(ctrlunit.Out.Functs))
	c.alu.SetPin(alu.In.Operand1, c.reg.GetPin(reg.Out.RS1Data))
	c.alu.SetPin(alu.In.Operand2, c.reg.GetPin(reg.Out.RS2Data))
	c.reg.SetPin(reg.In.Data, c.alu.GetPin(alu.Out.Result))

	return nil
}

func (c *Core) startComponents() error {
	for _, comp := range []device.Type{c.ctrl, c.reg, c.alu} {
		if err := comp.Run(); err != nil {
			return err
		}
	}

	return nil
}
