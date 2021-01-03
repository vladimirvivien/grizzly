package ctrlunit

import (
	"fmt"
	"log"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
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

		var (
			functs, rd, funct3, rs1, rs2, imm          uint32
			aluOp, aluSrc, memRen, memWen, wbSel, werf uint32
			memEnableCtrl                              datapath.Packet
		)

		for range c.clk.Ticks() {

			select {
			case inst := <-insts:
				opcode := isa.GetInstOpcode(inst)

				switch opcode {
				// register integer instructions
				case isa.Opcodes.R:
					fields := integer.Decode(inst)

					functs = fields.Functs()
					rd = fields.Rd
					funct3 = fields.Funct3
					rs1 = fields.Rs1
					rs2 = fields.Rs2
					imm = 0
					// controls
					aluOp = encodeAluOp(functs)
					aluSrc = 0
					memRen = 0
					wbSel = 0
					werf = 1

					memEnableCtrl = datapath.Packet{memRen, c.memRenOut}

				// register immediate
				case isa.Opcodes.RI:
					fields := integer.DecodeImm(inst)

					switch fields.Funct3 {
					case 0b001, 0b101:
						imm = fields.Shift
					default:
						imm = fields.Imm
					}

					functs = fields.Functs()
					rd = fields.Rd
					funct3 = fields.Funct3
					rs1 = fields.Rs1
					rs2 = 0
					// controls
					aluOp = encodeAluOp(functs)
					aluSrc = 1
					memRen = 0
					wbSel = 0
					werf = 1

					// mem read enable
					memEnableCtrl = datapath.Packet{memRen, c.memRenOut}

				// load instruction
				case isa.Opcodes.L:
					fields := load.Decode(inst)

					functs = 0
					rd = fields.Rd
					funct3 = fields.Funct3
					rs1 = fields.Rs1
					rs2 = 0
					imm = fields.Imm
					// controls
					aluOp = encodeAluOp(alu.Ops.Add)
					aluSrc = 1
					memRen = 1
					wbSel = 1
					werf = 1

					// mem read enable
					memEnableCtrl = datapath.Packet{memRen, c.memRenOut}

				// Store instruction
				case isa.Opcodes.S:
					fields := store.Decode(inst)

					functs = 0
					rd = 0
					funct3 = fields.Funct3
					rs1 = fields.Rs1
					rs2 = fields.Rs2
					imm = fields.Imm
					// controls
					aluOp = encodeAluOp(alu.Ops.Add)
					aluSrc = 1
					memWen = 1
					wbSel = 0
					werf = 0

					// mem write enable
					memEnableCtrl = datapath.Packet{memWen, c.memWenOut}
				default:
					panic(fmt.Sprintf("unsupported opcode: %07b", opcode))
				}

				// send to components
				datapath.Send(
					datapath.Packet{aluOp, c.aluOpOut},

					// reg-alu; source select
					datapath.Packet{rs1, c.rs1Out},
					datapath.Packet{rs2, c.rs2Out},
					datapath.Packet{imm, c.immOut},
					datapath.Packet{aluSrc, c.aluSrcOut},

					// memory
					datapath.Packet{funct3, c.memOpOut},
					memEnableCtrl,

					// alu-mem write back mux
					datapath.Packet{wbSel, c.wbSelOut},

					// reg writeback
					datapath.Packet{werf, c.werfOut},
					datapath.Packet{rd, c.rdOut},
				)
			}
		}
	}()
	return nil
}
