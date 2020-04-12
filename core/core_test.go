package core

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

func TestCore(t *testing.T) {
	tests := []struct {
		name      string
		core      func(*testing.T) *Core
		instructs func(*testing.T) device.WiresOut
		eval      func(*testing.T, *Core, chan struct{})
	}{
		{
			name: "single addition",
			core: func(t *testing.T) *Core {
				cor := newCore()
				regfile := cor.reg.(*reg.RegisterFile)
				regfile.SideLoad(0b00001, 10)
				regfile.SideLoad(0b00010, 7)
				return cor
			},
			instructs: func(t *testing.T) device.WiresOut {
				datapath := make(device.Wires)
				go func() {
					datapath <- 0b0000000_00010_00001_000_00101_0110011 // add reg[5], reg[1], reg[2]
				}()
				return datapath
			},
			eval: func(t *testing.T, cor *Core, waiter chan struct{}) {
				defer close(waiter)
				regfile := cor.reg.(*reg.RegisterFile)
				go func() {
					rd := regfile.Probe(0b00101)
					if rd != 17 {
						t.Fatalf("unexpected add operation result %b", rd)
					}
				}()

			},
		},
		{
			name: "multipe additions",
			core: func(t *testing.T) *Core {
				cor := newCore()
				regfile := cor.reg.(*reg.RegisterFile)
				regfile.SideLoad(0b00001, 4)
				regfile.SideLoad(0b00010, 7)
				regfile.SideLoad(0b01001, 12)
				return cor
			},
			instructs: func(t *testing.T) device.WiresOut {
				datapath := make(device.Wires)
				go func() {
					datapath <- 0b0000000_00010_00001_000_00101_0110011 // add reg[5], reg[1], reg[2]
					datapath <- 0b0000000_01001_00101_000_00011_0110011 // add reg[3], reg[5], reg[9]
				}()
				return datapath
			},
			eval: func(t *testing.T, cor *Core, waiter chan struct{}) {
				defer close(waiter)
				regfile := cor.reg.(*reg.RegisterFile)
				rd5 := regfile.Probe(0b00101)
				if rd5 != 11 {
					t.Fatalf("unexpected result for add: reg[5], reg[1], reg[2]: %b", rd5)
				}
				rd3 := regfile.Probe(0b00011)
				if rd3 != 23 {
					t.Fatalf("Unexpected result for add: reg[3], reg[5], reg[9] : %b", rd3)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})

			cor := test.core(t)
			if cor == nil {
				cor = newCore()
			}
			cor.SetPin(In.Insts, test.instructs(t))

			if err := cor.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				test.eval(t, cor, wait)
			}()

			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
