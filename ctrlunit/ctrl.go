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
		ALUOp  device.PinLabel
		ALUSrc device.PinLabel
		Werf   device.PinLabel
	}{
		RS1:    "ctrlunit.rs1.out",
		RS2:    "ctrlunit.rs2.out",
		RD:     "ctrlunit.rd.out",
		Imm:    "ctrlunit.immOut.out",
		ALUOp:  "ctrlunit.aluop.out",
		ALUSrc: "ctrlunit.alusrc.out",
		Werf:   "ctrlunit.werfOut.out",
	}
)

// Controller encodes the logic for the control unit
// It decodes the instruction and orchestrate the operation
// on the data using ALU, register file, etc.

// Data order: data is output on the following sequence:
// data path: rdOut, rs1, rs2,
// control path: aluOp, and werfOut
// If not read in that order, races will be created.
type Controller struct {
	*device.Base
	rdOut     device.Wires // regfile data address
	rs1Out    device.Wires // regfile select addr 1
	rs2Out    device.Wires // regfile select addr 2
	immOut    device.Wires // immediate value
	aluOpOut  device.Wires // ALU operation
	aluSrcOut device.Wires // ALU source mux selector
	werfOut   device.Wires // regfile write enable file
}

// New creates a new *Controller
func New() device.Type {
	return newCtrl()
}

func newCtrl() *Controller {
	c := &Controller{
		Base:      device.NewBase(),
		rdOut:     device.MakeWires(),
		rs1Out:    device.MakeWires(),
		rs2Out:    device.MakeWires(),
		immOut:    device.MakeWires(),
		aluOpOut:  device.MakeWires(),
		aluSrcOut: device.MakeWires(),
		werfOut:   device.MakeWires(),
	}
	c.SetPin(Out.RD, c.rdOut)
	c.SetPin(Out.RS1, c.rs1Out)
	c.SetPin(Out.RS2, c.rs2Out)
	c.SetPin(Out.Imm, c.immOut)
	c.SetPin(Out.ALUOp, c.aluOpOut)
	c.SetPin(Out.ALUSrc, c.aluSrcOut)
	c.SetPin(Out.Werf, c.werfOut)

	return c
}

func (c *Controller) Run() error {
	go func() {
		defer func() {
			close(c.rdOut)
			close(c.rs1Out)
			close(c.rs2Out)
			close(c.aluOpOut)
			close(c.aluSrcOut)
			close(c.werfOut)
		}()

		insts := c.GetPin(In.Insts)

		for {
			select {
			case inst := <-insts:
				opcode := inst & 0x7F
				switch opcode {
				case isa.Opcodes.R:
					fields := decodeR(inst)

					// controls
					c.werfOut <- 1
					c.aluOpOut <- encodeAluOp(fields.Functs())
					c.aluSrcOut <- 0

					// data path sent in order rs1,rs2,rd
					c.rs1Out <- fields.Rs1
					c.rs2Out <- fields.Rs2
					c.rdOut <- fields.Rd

				case isa.Opcodes.RI:
					fields := decodeRI(inst)

					// controls
					c.werfOut <- 1
					c.aluOpOut <- encodeAluOp(fields.Functs())
					c.aluSrcOut <- 1

					// data path order: rs1, imm, rd
					c.rs1Out <- fields.Rs1
					switch fields.Funct3 {
					case 0b001, 0b101:
						c.immOut <- fields.Shift
					default:
						c.immOut <- fields.Imm
					}
					c.rdOut <- fields.Rd

				default:
					panic(fmt.Sprintf("unsupported opcode: %0b", opcode))
				}
			}
		}
	}()
	return nil
}
