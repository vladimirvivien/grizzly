package ctrlunit

import (
	"fmt"
	"log"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
	"github.com/vladimirvivien/grizzly/isa/load"
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
		MemRen device.PinLabel
		MemWen device.PinLabel
		MemOp  device.PinLabel
		WBSel  device.PinLabel
	}{
		RS1:    "ctrlunit.rs1.out",
		RS2:    "ctrlunit.rs2.out",
		RD:     "ctrlunit.rd.out",
		Imm:    "ctrlunit.imm.out",
		ALUOp:  "ctrlunit.aluop.out",
		ALUSrc: "ctrlunit.alusrc.out",
		Werf:   "ctrlunit.werf.out",
		MemRen: "ctrlunit.memreadenable.out",
		MemWen: "ctrlunit.memwriteenable.out",
		MemOp:  "ctrlunit.memoperation.out",
		WBSel:  "ctrlunit.wbsel.out",
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
	memRenOut datapath.Wires // memory read enable
	memWenOut datapath.Wires // memory write enable
	memOpOut  datapath.Wires // memory operation control
	wbSelOut  datapath.Wires // register write-back selector
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
		memRenOut: datapath.MakeWires(),
		memWenOut: datapath.MakeWires(),
		memOpOut:  datapath.MakeWires(),
		wbSelOut:  datapath.MakeWires(),
	}
	c.SetPin(Out.RD, c.rdOut)
	c.SetPin(Out.RS1, c.rs1Out)
	c.SetPin(Out.RS2, c.rs2Out)
	c.SetPin(Out.Imm, c.immOut)
	c.SetPin(Out.ALUOp, c.aluOpOut)
	c.SetPin(Out.ALUSrc, c.aluSrcOut)
	c.SetPin(Out.Werf, c.werfOut)
	c.SetPin(Out.MemRen, c.memRenOut)
	c.SetPin(Out.MemWen, c.memWenOut)
	c.SetPin(Out.MemOp, c.memOpOut)
	c.SetPin(Out.WBSel, c.wbSelOut)

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
			close(c.memRenOut)
			close(c.memWenOut)
			close(c.memOpOut)
			close(c.wbSelOut)
		}()

		for range c.clk.Ticks() {

			select {
			case inst := <-insts:
				opcode := isa.GetInstOpcode(inst)

				switch opcode {

				// register integer instructions
				case isa.Opcodes.R:
					fields := integer.Decode(inst)

					// send register file ctrl and addrs
					datapath.Send(
						datapath.Packet{encodeAluOp(fields.Functs()), c.aluOpOut},

						// reg-alu; operand select
						datapath.Packet{fields.Rs1, c.rs1Out},
						datapath.Packet{fields.Rs2, c.rs2Out},
						datapath.Packet{0, c.immOut},
						datapath.Packet{0, c.aluSrcOut},

						// memory
						datapath.Packet{fields.Funct3, c.memOpOut},
						datapath.Packet{0, c.memRenOut},

						// write back mux
						datapath.Packet{0, c.wbSelOut},

						// alu-mem reg writeback
						datapath.Packet{1, c.werfOut},
						datapath.Packet{fields.Rd, c.rdOut},
					)

				// register immediate
				case isa.Opcodes.RI:
					fields := integer.DecodeImm(inst)

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

						// memory
						datapath.Packet{fields.Funct3, c.memOpOut},
						datapath.Packet{0, c.memRenOut},

						// alu-mem write back mux
						datapath.Packet{0, c.wbSelOut},

						// alu-mem; reg write back
						datapath.Packet{1, c.werfOut},
						datapath.Packet{fields.Rd, c.rdOut},
					)

				// load instruction
				case isa.Opcodes.L:
					fields := load.Decode(inst)
					datapath.Send(
						datapath.Packet{encodeAluOp(fields.Funct3), c.aluOpOut},

						// reg-alu; source select
						datapath.Packet{fields.Rs1, c.rs1Out},
						datapath.Packet{0, c.rs2Out},
						datapath.Packet{fields.Imm, c.immOut},
						datapath.Packet{1, c.aluSrcOut},

						// memory
						datapath.Packet{fields.Funct3, c.memOpOut},
						datapath.Packet{1, c.memRenOut},
						datapath.Packet{1, c.wbSelOut},

						// mem-reg
						datapath.Packet{1, c.werfOut},
						datapath.Packet{fields.Rd, c.rdOut},
					)
				default:
					panic(fmt.Sprintf("unsupported opcode: %07b", opcode))
				}

			}
		}
	}()
	return nil
}
