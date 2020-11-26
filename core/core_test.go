package core

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/reg"
)

func TestCore(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*testing.T, chan struct{}) *Core
		regMap map[uint32]datapath.Word
	}{

		{
			name: "multiple R and RIs",
			setup: func(t *testing.T, doneSignal chan struct{}) *Core {
				cor := newCore()
				regfile := cor.reg.(*reg.RegisterFile)
				regfile.SideLoad(1, 4)
				regfile.SideLoad(2, 2)

				insts := datapath.MakeWires()
				go func() {
					insts <- 0b000000000010_00001_000_00101_0010011  // addi reg[5] <= 2, reg[1]; reg[5]=6
					insts <- 0b0000000_00010_00101_000_00011_0110011 // add  reg[3] <= reg[5], reg[2]; reg[3]=8
					insts <- 0b0000000_00001_00011_001_00110_0010011 // slli reg[6] <= 1, reg[3]; reg[6]=16
					close(doneSignal)
				}()
				cor.SetPin(In.Insts, insts)
				return cor
			},

			regMap: map[uint32]datapath.Word{
				0b00101: 6,
				0b00011: 8,
				0b00110: 16,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			waiter := make(chan struct{})
			cor := test.setup(t, waiter)

			if err := cor.Run(); err != nil {
				t.Fatal(err)
			}

			select {
			case <-waiter:
				t.Log("waiting... before evaluation")
				time.Sleep(1000 * time.Millisecond)
				regfile := cor.reg.(*reg.RegisterFile)
				for k, expected := range test.regMap {
					probed := regfile.Probe(k)
					if probed != expected {
						t.Errorf("unexpected register value: reg[%05b]=%032b; expecting %032b", k, probed, expected)
					}
				}
			case <-time.After(5000 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
