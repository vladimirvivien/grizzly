package core

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

func TestCore(t *testing.T) {
	tests := []struct {
		name      string
		core      func(*testing.T) *Core
		instructs func(*testing.T) device.WiresOut
		eval      func(*testing.T, *Core)
		probeFor uint32
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

				insts := device.MakeWires()
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

			// fan out wire from alu output
			// connect one output back to register data
			// use other output for probing.
			pins := device.Fanout(cor.alu.GetPin(alu.Out.Result), 1)
			cor.reg.SetPin(reg.In.Data, pins[0])

			// probe alu output
			go func(){
				//t.Log(<-pins[1])
				//t.Log(<-pins[1])
				//t.Log(<-pins[1])
				time.Sleep(5*time.Second)
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
			case <-time.After(200000 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
