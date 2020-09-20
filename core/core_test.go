package core

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/ctrlunit"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

func TestCore(t *testing.T) {
	tests := []struct {
		name      string
		core      func(*testing.T) (*Core, chan struct{})
		instructs func(*testing.T) device.WiresOut
		eval      func(*testing.T, *Core)
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
			core: func(t *testing.T) (*Core, chan struct{}) {
				waiter := make(chan struct{})
				regfile := reg.New().(*reg.RegisterFile)
				regfile.Print()
				regfile.SideLoad(0b00001, 4)
				regfile.SideLoad(0b00010, 2)
				regfile.Print()

				ctrl := ctrlunit.New()
				disableCtrl := device.MakeWires()
				ctrl.SetPin(ctrlunit.In.Disable, disableCtrl)

				cor := newCore()
				cor.reg = regfile
				cor.ctrl = ctrl

				insts := make(device.Wires)
				go func() {
					insts <- 0b000000000010_00001_000_00101_0010011  // addi reg[5] <= 2, reg[1]
					insts <- 0b0000000_00010_00101_000_00011_0110011 // add  reg[3] <= reg[5], reg[2]
					insts <- 0b0000000_00001_00011_001_00110_0010011 // slli reg[6] <= 1, reg[3]
					close(waiter)
				}()
				cor.SetPin(In.Insts, insts)
				return cor, waiter
			},

			eval: func(t *testing.T, cor *Core) {
				//time.Sleep(4000 * time.Millisecond)
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
			cor, waiter := test.core(t)

			if err := cor.Run(); err != nil {
				t.Fatal(err)
			}

			select {
			case <-waiter:
				test.eval(t, cor)
			case <-time.After(5000 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
