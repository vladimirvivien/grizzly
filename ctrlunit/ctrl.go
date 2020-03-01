package ctrlunit

import (
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	In = struct {
		Insts device.PinLabel
	}{
		Insts: "ctrlunit.instructions.in",
	}

	Out = struct {
		RS1    device.PinLabel
		RS2    device.PinLabel
		RD     device.PinLabel
		Functs device.PinLabel
	}{
		RS1:    "ctrlunit.rs1.out",
		RS2:    "ctrlunit.rs2.out",
		RD:     "ctrlunit.rd.out",
		Functs: "ctrlunit.functs.out",
	}
)

type Controller struct {
	pins      device.Pins
	rs1Out    device.Wires
	rs2Out    device.Wires
	rdOut     device.Wires
	functsOut device.Wires
}

func New() device.Type {
	return newCtrl()
}

func newCtrl() *Controller {
	c := &Controller{
		pins:      make(device.Pins),
		functsOut: device.MakeWires(),
		rs1Out:    device.MakeWires(),
		rs2Out:    device.MakeWires(),
		rdOut:     device.MakeWires(),
	}
	c.pins = device.Pins{
		Out.Functs: c.functsOut,
		Out.RS1:    c.rs1Out,
		Out.RS2:    c.rs2Out,
		Out.RD:     c.rdOut,
	}
	return c
}

func (c *Controller) GetPins() device.Pins {
	return c.pins
}

func (c *Controller) GetPin(label device.PinLabel) device.Pin {
	return c.pins[label]
}

func (c *Controller) SetPin(label device.PinLabel, pin device.Pin) {
	c.pins[label] = pin
}

func (c *Controller) Run() error {
	go func() {
		defer func() {
			close(c.functsOut)
			close(c.rs1Out)
			close(c.rs2Out)
			close(c.rdOut)
		}()

		for {
			inst := <-c.GetPin(In.Insts)
			opcode := inst & 0x7F

			switch opcode {
			case isa.Opcodes.R:
				fields := decodeR(inst)
				c.functsOut <- fields.Functs()
				c.rs1Out <- fields.Rs1
				c.rs2Out <- fields.Rs2
				c.rdOut <- fields.Rd
			}
		}
	}()
	return nil
}
