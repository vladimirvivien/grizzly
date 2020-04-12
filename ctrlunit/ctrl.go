package ctrlunit

import (
	"fmt"

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
		Werf   device.PinLabel
	}{
		RS1:    "ctrlunit.rs1.out",
		RS2:    "ctrlunit.rs2.out",
		RD:     "ctrlunit.rd.out",
		Functs: "ctrlunit.functs.out",
		Werf:   "ctrlunit.werf.out",
	}
)

type Controller struct {
	*device.Base
	rs1Out    device.Wires
	rs2Out    device.Wires
	rdOut     device.Wires
	werf      device.Wires
	functsOut device.Wires
}

func New() device.Type {
	return newCtrl()
}

func newCtrl() *Controller {
	c := &Controller{
		Base:      device.NewBase(),
		functsOut: device.MakeWires(),
		rs1Out:    device.MakeWires(),
		rs2Out:    device.MakeWires(),
		rdOut:     device.MakeWires(),
		werf:      device.MakeWires(),
	}
	c.SetPin(Out.Functs, c.functsOut)
	c.SetPin(Out.RS1, c.rs1Out)
	c.SetPin(Out.RS2, c.rs2Out)
	c.SetPin(Out.RD, c.rdOut)
	c.SetPin(Out.Werf, c.werf)

	return c
}

func (c *Controller) Run() error {
	fmt.Println("Controller started")
	go func() {
		defer func() {
			close(c.functsOut)
			close(c.rs1Out)
			close(c.rs2Out)
			close(c.rdOut)
			close(c.werf)
		}()

		for {
			inst := <-c.GetPin(In.Insts)
			opcode := inst & 0x7F

			switch opcode {
			case isa.Opcodes.R:
				// R-format instructions
				// decodes instructions
				// read operands RS1, RS2
				// Exec ALU operation
				// write back result RD
				fields := decodeR(inst)
				c.rs1Out <- fields.Rs1
				c.rs2Out <- fields.Rs2
				c.functsOut <- fields.Functs()
				c.werf <- 1 // TODO change to bit type
				c.rdOut <- fields.Rd

			default:
				panic(fmt.Sprintf("unsupported opcode: %0b", opcode))
			}
		}
	}()
	return nil
}
