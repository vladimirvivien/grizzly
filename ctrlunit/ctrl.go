package ctrlunit

import (
	"fmt"
	"log"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
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

// Order:
// The order in which data is sent to component matters
// as it may cause deadlock if subcomponents process data
// out of order.
type Controller struct {
	*device.Base
	rdOut     datapath.Wires // regfile data address
	rs1Out    datapath.Wires // regfile select addr 1
	rs2Out    datapath.Wires // regfile select addr 2
	immOut    datapath.Wires // immediate value
	aluOpOut  datapath.Wires // ALU operation
	aluSrcOut datapath.Wires // ALU source mux selector
	werfOut   datapath.Wires // regfile write enable file
	clk       clock.Clock
}

// New creates a new *Controller
func New() device.ClockedType {
	return newCtrl()
}

func newCtrl() *Controller {
	c := &Controller{
		Base:      device.NewBase(),
		rdOut:     datapath.MakeWires(),
		rs1Out:    datapath.MakeWires(),
		rs2Out:    datapath.MakeWires(),
		immOut:    datapath.MakeWires(),
		aluOpOut:  datapath.MakeWires(),
		aluSrcOut: datapath.MakeWires(),
		werfOut:   datapath.MakeWires(),
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

func (c *Controller) SetClock(clk clock.Clock) {
	c.clk = clk
}

func (c *Controller) Run() error {
	log.Println("controller: starting...")
	if c.clk == nil {
		panic("ctrl: missing clock")
	}

	insts := c.GetPin(In.Insts)
	if insts == nil {
		panic("ctrl: missing instructions")
	}

	go func() {
		defer func() {
			close(c.rdOut)
			close(c.rs1Out)
			close(c.rs2Out)
			close(c.aluOpOut)
			close(c.aluSrcOut)
			close(c.werfOut)
		}()

		for range c.clk.Ticks() {

			select {
			case inst := <-insts:
				opcode := inst & 0x7F

				switch opcode {
				case isa.Opcodes.R:
					fields := decodeR(inst)

					// send register file ctrl and addrs
					datapath.Send(
						datapath.Packet{encodeAluOp(fields.Functs()), c.aluOpOut},

						// reg-alu; operand select
						datapath.Packet{fields.Rs1, c.rs1Out},
						datapath.Packet{fields.Rs2, c.rs2Out},
						datapath.Packet{0, c.immOut},
						datapath.Packet{0, c.aluSrcOut},

						// alu-reg; data store
						datapath.Packet{1, c.werfOut},
						datapath.Packet{fields.Rd, c.rdOut},
					)

				case isa.Opcodes.RI:
					fields := decodeRI(inst)

					var imm uint32
					switch fields.Funct3 {
					case 0b001, 0b101:
						imm = fields.Shift
					default:
						imm = fields.Imm
					}

					datapath.Send(
						datapath.Packet{encodeAluOp(fields.Functs()), c.aluOpOut},

						// reg-alu; source select
						datapath.Packet{fields.Rs1, c.rs1Out},
						datapath.Packet{0, c.rs2Out},
						datapath.Packet{imm, c.immOut},
						datapath.Packet{1, c.aluSrcOut},

						// alu-reg; data store
						datapath.Packet{1, c.werfOut},
						datapath.Packet{fields.Rd, c.rdOut},
					)
				default:
					panic(fmt.Sprintf("unsupported opcode: %0b", opcode))
				}
			}
		}
	}()
	return nil
}
