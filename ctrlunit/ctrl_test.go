package ctrlunit

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

func TestCtrl_R(t *testing.T) {
	instructions := datapath.MakeWires()
	ctrl := newCtrl()
	ctrl.SetPin(In.Insts, instructions)

	tests := []struct {
		name string
		inst func() isa.Inst
		eval func(aluOp uint32, rs1 uint32, rs2 uint32, werf uint32, rd uint32)
	}{
		{
			name: "R format",
			inst: func() isa.Inst { return 0b0000000_00010_00001_000_00101_0110011 },
			eval: func(aluOp, rs1, rs2, werf, rd uint32) {
				if aluOp != alu.Ops.Add {
					t.Errorf("Unexpected Operation value: %b", aluOp)
				}
				if rs1 != 0b00001 {
					t.Errorf("Unexpected rs1 value: %b", rs1)
				}
				if rs2 != 0b00010 {
					t.Errorf("Unexpected rs2 value: %b", rs2)
				}
				if rd != 0b00101 {
					t.Errorf("Unexpected rd value: %b", rd)
				}
				if werf != 1 {
					t.Error("Unexpedted WERF value")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})

			go func() {
				instructions <- test.inst()
			}()

			if err := ctrl.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				defer close(wait)
				out := datapath.Collect(
					ctrl.GetPin(Out.ALUOp),
					ctrl.GetPin(Out.RS1),
					ctrl.GetPin(Out.RS2),
					ctrl.GetPin(Out.Imm),
					ctrl.GetPin(Out.ALUSrc),
					ctrl.GetPin(Out.Werf),
					ctrl.GetPin(Out.RD),
				)
				test.eval(out[0], out[1], out[2],out[5], out[6])
			}()

			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}

func TestCtrl_RI(t *testing.T) {
	instructions := datapath.MakeWires()
	ctrl := newCtrl()
	ctrl.SetPin(In.Insts, instructions)

	tests := []struct {
		name string
		inst func() isa.Inst
		eval func(aluOp uint32, rs1 uint32, imm uint32, werf uint32, rd uint32)
	}{
		{
			name: "RI format (addi)",
			inst: func() isa.Inst { return 0b000000000010_00001_000_00101_0010011 },
			eval: func(aluOp, rs1, imm, werf, rd uint32) {
				if aluOp != alu.Ops.Add {
					t.Errorf("Unexpected Operation value: %b", aluOp)
				}
				if rs1 != 0b00001 {
					t.Errorf("Unexpected rs1 value: %b", rs1)
				}
				if rd != 0b00101 {
					t.Errorf("Unexpected rd value: %b", rd)
				}
				if werf != 1 {
					t.Error("Unexpedted WERF value")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})

			go func() {
				instructions <- test.inst()
			}()

			if err := ctrl.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				defer close(wait)
				out := datapath.Collect(
					ctrl.GetPin(Out.ALUOp),
					ctrl.GetPin(Out.RS1),
					ctrl.GetPin(Out.RS2),
					ctrl.GetPin(Out.Imm),
					ctrl.GetPin(Out.ALUSrc),
					ctrl.GetPin(Out.Werf),
					ctrl.GetPin(Out.RD),
				)
				test.eval(out[0], out[1], out[3],out[5], out[6])
			}()

			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatalf("Control unit operation %s took too long", test.name)
			}
		})
	}
}
