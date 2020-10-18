package core

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/ctrlunit"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

// TestCore_CtrlReg
// Tests execution between controller and register file
func TestCore_CtrlReg(t *testing.T) {
	tests := []struct {
		name string
		prep func(*testing.T) (*Core, datapath.Wires)
		eval func(*testing.T, *Core, datapath.Wires)
	}{
		{
			name: "add",
			prep: func(t *testing.T, ) (*Core, datapath.Wires) {
				cor := newCore()
				insts := datapath.MakeWires()
				cor.ctrl.SetPin(ctrlunit.In.Insts, insts)
				regData := datapath.MakeWires()

				reg := cor.reg.(*reg.RegisterFile)
				reg.SideLoad(2, 4)
				reg.SideLoad(6, 12)
				reg.SideLoad(8, 16)

				go func() {
					insts <- 0b0000000_00110_00010_000_00101_0110011 // add  reg[5]  = reg[2]=4, reg[6]=12
					insts <- 0b000000000010_00101_000_00101_0010011  // addi reg[5]  = reg[5]=16, 2
					insts <- 0b0000000_01000_00110_000_00111_0110011 // add  reg[7]  = reg[6]=12, reg[8]=16
					insts <- 0b0000000_00010_01000_000_01010_0110011 // add  reg[10] = reg[8]=16, reg[2]=4
				}()

				return cor, regData
			},
			eval: func(t *testing.T, cor *Core, regData datapath.Wires) {
				ctrl := cor.ctrl.(*ctrlunit.Controller)
				aluOp := ctrl.GetPin(ctrlunit.Out.ALUOp)
				aluSrc := ctrl.GetPin(ctrlunit.Out.ALUSrc)
				imm := ctrl.GetPin(ctrlunit.Out.Imm)

				regfile := cor.reg.(*reg.RegisterFile)
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
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cor, regData := test.prep(t)

			// wire the register
			cor.reg.SetPin(reg.In.Werf, cor.ctrl.GetPin(ctrlunit.Out.Werf))
			cor.reg.SetPin(reg.In.RS1Addr, cor.ctrl.GetPin(ctrlunit.Out.RS1))
			cor.reg.SetPin(reg.In.RS2Addr, cor.ctrl.GetPin(ctrlunit.Out.RS2))
			cor.reg.SetPin(reg.In.RDAddr, cor.ctrl.GetPin(ctrlunit.Out.RD))
			cor.reg.SetPin(reg.In.Data, regData)

			if err := cor.reg.Run(); err != nil {
				t.Fatal(err)
			}
			if err := cor.ctrl.Run(); err != nil {
				t.Fatal(err)
			}

			test.eval(t, cor, regData)
		})
	}
}

func TestCore(t *testing.T) {
	tests := []struct {
		name      string
		core      func(*testing.T) *Core
		instructs func(*testing.T) device.Pin
		eval      func(*testing.T, *Core)
		probeFor  uint32
	}{
		//{
		//	name: "single R instruction",
		//	core: func(t *testing.T) *Core {
		//		cor := newCore()
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		regfile.SideLoad(0b00001, 4)
		//		regfile.SideLoad(0b00010, 7)
		//		return cor
		//	},
		//	instructs: func(t *testing.T) device.WiresOut {
		//		datapath := make(device.Wires)
		//		go func() {
		//			datapath <- 0b0000000_00010_00001_000_00101_0110011 // add  reg[5] <= reg[1], reg[2]
		//		}()
		//		return datapath
		//	},
		//	eval: func(t *testing.T, cor *Core, waiter chan struct{}) {
		//		defer close(waiter)
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		rd5 := regfile.Probe(0b00101)
		//		if rd5 != 11 {
		//			t.Fatalf("unexpected result for add: reg[5], reg[1], reg[2]: %b", rd5)
		//		}
		//	},
		//},
		//{
		//	name: "multiple R instructions",
		//	core: func(t *testing.T) *Core {
		//		cor := newCore()
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		regfile.SideLoad(0b00001, 4)
		//		regfile.SideLoad(0b00010, 7)
		//		regfile.SideLoad(0b01001, 12)
		//		return cor
		//	},
		//	instructs: func(t *testing.T) device.WiresOut {
		//		datapath := make(device.Wires)
		//		go func() {
		//			datapath <- 0b0000000_00010_00001_000_00101_0110011 // add  reg[5] <= reg[1], reg[2]
		//			datapath <- 0b0000000_01001_00101_000_00011_0110011 // add  reg[3] <= reg[5], reg[9]
		//		}()
		//		return datapath
		//	},
		//	eval: func(t *testing.T, cor *Core, waiter chan struct{}) {
		//		defer close(waiter)
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		rd5 := regfile.Probe(0b00101)
		//		if rd5 != 11 {
		//			t.Fatalf("unexpected result for add: reg[5], reg[1], reg[2]: %b", rd5)
		//		}
		//		rd3 := regfile.Probe(0b00011)
		//		if rd3 != 23 {
		//			t.Fatalf("Unexpected result for add: reg[3], reg[5], reg[9] : %b", rd3)
		//		}
		//	},
		//},
		//{
		//	name: "single RI add",
		//	core: func(t *testing.T) *Core {
		//		cor := newCore()
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		regfile.SideLoad(0b00001, 4)
		//		return cor
		//	},
		//	instructs: func(t *testing.T) device.WiresOut {
		//		datapath := make(device.Wires)
		//		go func() {
		//			datapath <- 0b000000000010_00001_000_00101_0010011 // addi reg[5] <= 2, reg[1]
		//		}()
		//		return datapath
		//	},
		//	eval: func(t *testing.T, cor *Core, waiter chan struct{}) {
		//		defer close(waiter)
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		rd5 := regfile.Probe(0b00101)
		//		if rd5 != 6 {
		//			t.Fatalf("unexpected result for addi: reg[5] <= 2, reg[1]: %d", rd5)
		//		}
		//	},
		//},
		//{
		//	name: "multiple RI",
		//	core: func(t *testing.T) *Core {
		//		cor := newCore()
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		regfile.SideLoad(0b00001, 4)
		//		return cor
		//	},
		//	instructs: func(t *testing.T) device.WiresOut {
		//		datapath := make(device.Wires)
		//		go func() {
		//			datapath <- 0b000000000010_00001_000_00101_0010011 // addi reg[5] <= 2, reg[1]
		//			datapath <- 0b0000000_0001_00101_001_00110_0010011 // slli reg[6] <= 1, reg[5]
		//		}()
		//		return datapath
		//	},
		//	eval: func(t *testing.T, cor *Core, waiter chan struct{}) {
		//		defer close(waiter)
		//		regfile := cor.reg.(*reg.RegisterFile)
		//		rd5 := regfile.Probe(0b00101)
		//		if rd5 != 6 {
		//			t.Fatalf("unexpected result for addi: reg[5] <= 2, reg[1]: %d", rd5)
		//		}
		//		rd6 := regfile.Probe(0b00110)
		//		if rd6 != 12 {
		//			t.Fatalf("unexpected result for slli: reg[6] <= 1, reg[5]: %d", rd5)
		//		}
		//	},
		//},
		{
			name: "multiple R and RIs",
			core: func(t *testing.T) *Core {
				regfile := reg.New().(*reg.RegisterFile)
				regfile.SideLoad(0b00001, 4)
				regfile.SideLoad(0b00010, 2)

				cor := newCore()
				cor.reg = regfile

				insts := datapath.MakeWires()
				go func() {
					insts <- 0b000000000010_00001_000_00101_0010011  // addi reg[5] <= 2, reg[1]; reg[5]=6
					insts <- 0b0000000_00010_00101_000_00011_0110011 // add  reg[3] <= reg[5], reg[2]; reg[3]=8
					insts <- 0b0000000_00001_00011_001_00110_0010011 // slli reg[6] <= 1, reg[3]; reg[6]=16
				}()
				cor.SetPin(In.Insts, insts)
				return cor
			},

			probeFor: 0b10000,

			eval: func(t *testing.T, cor *Core) {
				regfile := cor.reg.(*reg.RegisterFile)
				regfile.Print()

				//rd5 := regfile.Probe(0b00101)
				//if rd5 != 6 {
				//	t.Fatalf("unexpected result for addi: reg[5] <= 2, reg[1]: %d", rd5)
				//}
				//rd3 := regfile.Probe(0b00011)
				//if rd3 != 8 {
				//	t.Fatalf("unexpected result for add: reg[3] <= reg[5], reg[2]: %d", rd3)
				//}
				//rd6 := regfile.Probe(0b00110)
				//if rd6 != 16 {
				//	t.Fatalf("unexpected result for slli: reg[6] <= 1, reg[3]: %d", rd5)
				//}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			waiter := make(chan struct{})
			cor := test.core(t)
			regfile := cor.reg.(*reg.RegisterFile)

			// fan out wire from alu output
			// connect one output back to register data
			// use other output for probing.
			pins := device.Fanout(cor.alu.GetPin(alu.Out.Result), 1)
			cor.reg.SetPin(reg.In.Data, pins[0])

			// probe alu output
			go func() {
				//t.Log(<-pins[1])
				//t.Log(<-pins[1])
				//t.Log(<-pins[1])
				ticker := time.NewTicker(100 * time.Millisecond)
				for {
					select {
					case <-ticker.C:
						if regfile.Probe(6) == 0 {
							continue
						}
						break

					case <-time.After(1000 * time.Millisecond):
						t.Fatal("tired of waiting for probe value")
					}

				}
				close(waiter)
				//for{
				//	select{
				//	case  val := <- pins[1]:
				//		if val == test.probeFor{
				//			fmt.Println(val)
				//			close(waiter)
				//		}
				//	}
				//}
			}()

			if err := cor.Run(); err != nil {
				t.Fatal(err)
			}

			select {
			case <-waiter:
				t.Log("Eval()")
				test.eval(t, cor)
			case <-time.After(5000 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
