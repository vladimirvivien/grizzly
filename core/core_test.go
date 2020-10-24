package core

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/ctrlunit"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

// TestCore_CtrlReg
// Tests execution between controller and register file
func TestCore_CtrlReg(t *testing.T) {
	cor := newCore()

	insts := datapath.MakeWires()
	cor.ctrl.SetPin(ctrlunit.In.Insts, insts)
	regData := datapath.MakeWires()

	// connect controller and register
	cor.reg.SetPin(reg.In.Werf, cor.ctrl.GetPin(ctrlunit.Out.Werf))
	cor.reg.SetPin(reg.In.RS1Addr, cor.ctrl.GetPin(ctrlunit.Out.RS1))
	cor.reg.SetPin(reg.In.RS2Addr, cor.ctrl.GetPin(ctrlunit.Out.RS2))
	cor.reg.SetPin(reg.In.RDAddr, cor.ctrl.GetPin(ctrlunit.Out.RD))
	cor.reg.SetPin(reg.In.Data, regData)

	// prepare register
	regfile := cor.reg.(*reg.RegisterFile)
	regfile.SideLoad(2, 4)
	regfile.SideLoad(6, 12)
	regfile.SideLoad(8, 16)

	// instructions
	go func() {
		insts <- 0b0000000_00110_00010_000_00101_0110011 // add  reg[5]  = reg[2]=4, reg[6]=12
		insts <- 0b000000000010_00101_000_00101_0010011  // addi reg[5]  = reg[5]=16, 2
		insts <- 0b0000000_01000_00110_000_00111_0110011 // add  reg[7]  = reg[6]=12, reg[8]=16
		insts <- 0b0000000_00010_01000_000_01010_0110011 // add  reg[10] = reg[8]=16, reg[2]=4
	}()

	// start components
	if err := cor.reg.Run(); err != nil {
		t.Fatal(err)
	}
	if err := cor.ctrl.Run(); err != nil {
		t.Fatal(err)
	}

	ctrl := cor.ctrl.(*ctrlunit.Controller)
	aluOp := ctrl.GetPin(ctrlunit.Out.ALUOp)
	aluSrc := ctrl.GetPin(ctrlunit.Out.ALUSrc)
	imm := ctrl.GetPin(ctrlunit.Out.Imm)

	rs1Data := regfile.GetPin(reg.Out.RS1Data)
	rs2Data := regfile.GetPin(reg.Out.RS2Data)

	// collect from
	data := datapath.Collect(aluOp, aluSrc, rs1Data, rs2Data)
	op1, op2 := data[2], data[3]
	if op1 != 4 {
		t.Fatalf("Unexpected data from register line %d", op1)
	}
	if op2 != 12 {
		t.Fatalf("Unexpected data from register line %d", op2)
	}
	// register file data line must be provided after each R inst or deadlock will happen
	regData <- op1 + op2

	data = datapath.Collect(aluOp, aluSrc, rs1Data, rs2Data, imm)
	op1, op2, immOp := data[2], data[3], data[4]
	if op1 != 16 {
		t.Fatalf("Unexpected op1 data: %d", op1)
	}
	if op2 != 0 {
		t.Fatalf("reg op2 should be 0 in imm, bot got %d", op2)
	}
	if immOp != 2 {
		t.Fatalf("unexpected data for ctrl imm %d", immOp)
	}
	regData <- op1 + immOp

	data = datapath.Collect(aluOp, aluSrc, rs1Data, rs2Data)
	op1, op2 = data[2], data[3]
	if op1 != 12 {
		t.Fatalf("Unexpected data from register line %d", op1)
	}
	if op2 != 16 {
		t.Fatalf("Unexpected data from register line %d", op2)
	}
	regData <- op1 + op2

	data = datapath.Collect(aluOp, aluSrc, rs1Data, rs2Data)
	op1, op2 = data[2], data[3]
	if op1 != 16 {
		t.Fatalf("Unexpected data from register line %d", op1)
	}
	if op2 != 4 {
		t.Fatalf("Unexpected data from register line %d", op2)
	}
	regData <- op1 + op2

}

func TestCore(t *testing.T) {
	tests := []struct {
		name      string
		core      func(*testing.T, chan struct{}) *Core
		instructs func(*testing.T) device.Pin
		eval      func(*testing.T, *Core)
		probeFor  uint32
	}{

		{
			name: "multiple R and RIs",
			core: func(t *testing.T, stopSignal chan struct{}) *Core {
				cor := newCore()
				regfile := cor.reg.(*reg.RegisterFile)
				regfile.SideLoad(1, 4)
				regfile.SideLoad(2, 2)

				insts := datapath.MakeWires()
				go func() {
					insts <- 0b000000000010_00001_000_00101_0010011  // addi reg[5] <= 2, reg[1]; reg[5]=6
					insts <- 0b0000000_00010_00101_000_00011_0110011 // add  reg[3] <= reg[5], reg[2]; reg[3]=8
					insts <- 0b0000000_00001_00011_001_00110_0010011 // slli reg[6] <= 1, reg[3]; reg[6]=16
					close(stopSignal)
				}()
				cor.SetPin(In.Insts, insts)
				return cor
			},

			eval: func(t *testing.T, cor *Core) {
				t.Log("Evaluating...")
				regfile := cor.reg.(*reg.RegisterFile)
				// reg[1]
				if val := regfile.Probe(0b00001); val != 4 {
					t.Errorf("unexpected value: reg[%05b]= %032b", 1, val)
				}
				// reg[2]
				if val := regfile.Probe(0b00010); val != 2 {
					t.Errorf("unexpected value: reg[%05b]= %032b", 2, val)
				}
				// reg[5]
				if val := regfile.Probe(0b00101); val != 6 {
					t.Errorf("unexpected value: reg[%05b]= %032b", 5, val)
				}
				// reg[3]
				if val := regfile.Probe(0b00011); val != 8 {
					t.Errorf("unexpected value: reg[%05b]= %032b", 3, val)
				}
				// reg[6]
				if val := regfile.Probe(0b00011); val != 16 {
					t.Errorf("unexpected value: reg[%05b]= %032b", 6, val)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			waiter := make(chan struct{})
			cor := test.core(t, waiter)

			if err := cor.Run(); err != nil {
				t.Fatal(err)
			}

			select {
			case <-waiter:
				t.Log("waiting... before evaluation")
				time.Sleep(200 * time.Millisecond)
				test.eval(t, cor)
			case <-time.After(5000 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
