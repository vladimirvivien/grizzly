package core

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/device"
)

func TestCore(t *testing.T) {
	tests := []struct {
		name      string
		instructs func() device.WiresOut
		eval      func(*Core, chan struct{})
	}{
		{
			name: "single addition",
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
				//regfile := cor.reg.(*reg.RegisterFile)

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})

			cor := newCore()
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
