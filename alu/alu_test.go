package alu

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/inst"
)

func TestALUOps(t *testing.T) {
	tests := []struct {
		name      string
		datalines func() (<-chan uint32, <-chan uint32, <-chan uint32)
		eval      func(t *testing.T, alu *ALU)
	}{
		{
			name: "add",
			datalines: func() (<-chan uint32, <-chan uint32, <-chan uint32) {
				d1wires := make(chan uint32)
				d2wires := make(chan uint32)
				opwires := make(chan uint32)
				go func() {
					defer func() {
						close(d1wires)
						close(d2wires)
					}()
					// add
					d1wires <- 0x2
					d2wires <- 0x2
					opwires <- inst.Add.Funct

					// sub
					d1wires <- 0x7
					d2wires <- 0x3
					opwires <- inst.Sub.Funct
				}()
				return d1wires, d2wires, opwires
			},
			eval: func(t *testing.T, alu *ALU) {
				result := <-alu.DataOut()
				if result != 0x4 {
					t.Error("unexpected result from alu addition 0x2 + 0x2:", result)
				}

				result = <-alu.DataOut()
				if result != 0x4 {
					t.Error("unexpected result from alu addition 0x2 + 0x2:", result)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})
			alu := newAlu()

			// wire ports
			d1wire, d2wire, opwire := test.datalines()
			alu.Data1In(d1wire)
			alu.Data2In(d2wire)
			alu.FunctIn(opwire)

			if err := alu.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				defer close(wait)
				test.eval(t, alu)
			}()

			// detect stuck path
			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatal("Register operation took too long to comlete...")
			}
		})
	}
}
