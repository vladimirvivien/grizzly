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
		Imm    device.PinLabel
		Functs device.PinLabel
		Werf   device.PinLabel
	}{
		RS1:    "ctrlunit.rs1.out",
		RS2:    "ctrlunit.rs2.out",
		RD:     "ctrlunit.rd.out",
		Imm:    "ctrlunit.imm.out",
		Functs: "ctrlunit.functs.out",
		Werf:   "ctrlunit.werf.out",
	}
)

// Controller encodes the logic for the control unit
// It decodes the instruction and orchestrate the operation
// on the data using ALU, register file, etc.

// Data order: data is output on the following sequence:
// data path: rdOut, rs1, rs2,
// control path: pathfuncts, and werf
// If not read in that order, races will be created.
type Controller struct {
	*device.Base
	rdOut     device.Wires
	rs1Out    device.Wires
	rs2Out    device.Wires
	imm       device.Wires
	functsOut device.Wires
	werf      device.Wires
}

// New creates a new *Controler
func New() device.Type {
	return newCtrl()
}

func newCtrl() *Controller {
	c := &Controller{
		Base:      device.NewBase(),
		rdOut:     device.MakeWires(),
		rs1Out:    device.MakeWires(),
		rs2Out:    device.MakeWires(),
		imm:       device.MakeWires(),
		functsOut: device.MakeWires(),
		werf:      device.MakeWires(),
	}
	c.SetPin(Out.RD, c.rdOut)
	c.SetPin(Out.RS1, c.rs1Out)
	c.SetPin(Out.RS2, c.rs2Out)
	c.SetPin(Out.Imm, c.imm)
	c.SetPin(Out.Functs, c.functsOut)
	c.SetPin(Out.Werf, c.werf)

	return c
}

func (c *Controller) Run() error {
	fmt.Println("Controller started")
	go func() {
		defer func() {
			close(c.rdOut)
			close(c.rs1Out)
			close(c.rs2Out)
			close(c.functsOut)
			close(c.werf)
		}()

		for {
			inst := <-c.GetPin(In.Insts)
			opcode := inst & 0x7F

			switch opcode {
			case isa.Opcodes.R:
				// R-format:
				fields := decodeR(inst)
				go func() {
					c.rdOut <- fields.Rd
					c.rs1Out <- fields.Rs1
					c.rs2Out <- fields.Rs2
				}()
				go func() {
					c.werf <- 1
				}()
				go func() {
					c.functsOut <- fields.Functs()
				}()
			case isa.Opcodes.RI:
				// RI-format (register immediate):
				fields := decodeRI(inst)

				go func() {
					c.rdOut <- fields.Rd
					c.rs1Out <- fields.Rs1

					// select Imm value or shift amout
					switch fields.Funct3 {
					case 0b001, 0b101:
						c.imm <- fields.Shift
					default:
						c.imm <- fields.Imm
					}

				}()

				go func() {
					c.functsOut <- isa.Functs(fields.Funct7, fields.Funct3)
				}()

				go func() {
					c.werf <- 1
				}()
			default:
				panic(fmt.Sprintf("unsupported opcode: %0b", opcode))
			}
		}
	}()
	return nil
}
