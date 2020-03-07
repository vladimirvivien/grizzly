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
		core      func() *Core
		instructs func() device.WiresOut
		eval      func(*Core, chan struct{})
	}{
		{
			name: "single addition",
			core: func() *Core {
				cor := newCore()
				regfile := cor.reg.(*reg.RegisterFile)
				regfile.SideLoad(0b00001, 10)
				regfile.SideLoad(0b00010, 7)
				return cor
			},
			instructs: func() device.WiresOut {
				datapath := make(device.Wires)
				go func() {
					datapath <- 0b0000000_00010_00001_000_00101_0110011 // add
				}()
				return datapath
			},
			eval: func(cor *Core, waiter chan struct{}) {
				defer close(waiter)

				// TODO extend reg.RegisterFile with methods
				// for inspections during testsw
				regfile := cor.reg.(*reg.RegisterFile)
				rd := regfile.Probe(0b00101)
				if rd != 17 {
					t.Fatalf("unexpected add operation result %b", rd)
				}

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})

			cor := test.core()
			if cor == nil {

				cor = newCore()
			}
			cor.SetPin(In.Insts, test.instructs())

			if err := cor.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				test.eval(cor, wait)
			}()

			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
