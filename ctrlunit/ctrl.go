package ctrlunit

import (
	"fmt"
	"time"

	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	In = struct {
		Insts   device.PinLabel
		Disable device.PinLabel
	}{
		Insts:   "ctrlunit.instructions.in",
		Disable: "ctrlunit.disable.in",
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

		for {
			select {
			case delay := <-c.GetPin(In.Disable):
				// this is for TEST-ONLY
				// It stalls for delay (micro sec) time to allow
				// inflight instructions to complete.
				// Then close all channels.
				time.Sleep(time.Duration(delay) * time.Microsecond)
				fmt.Printf("Disabling Ctrl after %d microsec ", delay)
				return
			case inst, valid := <-c.GetPin(In.Insts):
				if !valid {
					// this should never happen
					panic("controller: instruction channel is closed")
				}
				opcode := inst & 0x7F

				switch opcode {
				case isa.Opcodes.R:
					// R-format:
					fields := decodeR(inst)

					go func() {
						c.rs1Out <- fields.Rs1
					}()

					go func() {
						c.rs2Out <- fields.Rs2
					}()

					go func() {
						c.rdOut <- fields.Rd
					}()

					go func() {
						c.werfOut <- 1
					}()

					go func() {
						c.aluOpOut <- encodeAluOp(fields.Functs())
					}()

					go func() {
						c.aluSrcOut <- 0
					}()
				case isa.Opcodes.RI:
					// RI-format (register immediate):
					fields := decodeRI(inst)

					go func() {
						c.rdOut <- fields.Rd
					}()

					go func() {
						c.rs1Out <- fields.Rs1
					}()

					go func() {
						// select Imm value or shift amount
						switch fields.Funct3 {
						case 0b001, 0b101:
							c.immOut <- fields.Shift
						default:
							c.immOut <- fields.Imm
						}
					}()

					go func() {
						c.aluOpOut <- encodeAluOp(fields.Functs())
					}()

					go func() {
						c.aluSrcOut <- 1
					}()

					go func() {
						c.werfOut <- 1
					}()
				default:
					panic(fmt.Sprintf("unsupported opcode: %0b", opcode))
				}
			}
		}
	}()
	return nil
}
