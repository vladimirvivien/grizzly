//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package brancher

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/branch"
)

func TestBrancher_RV32(t *testing.T) {
	tests := []struct {
		name       string
		op         datapath.BranchOp
		wantTaken  bool
		wantOffset int32
	}{
		{
			name: "beq taken",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Beq.F3,
				RS1D:   42,
				RS2D:   42,
				Imm:    0b000000001010, // 10
			},
			wantTaken:  true,
			wantOffset: 20,
		},
		{
			name: "beq not taken",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Beq.F3,
				RS1D:   42,
				RS2D:   43,
				Imm:    0b000000001010,
			},
			wantTaken:  false,
			wantOffset: 4,
		},
		{
			name: "bne taken",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Bne.F3,
				RS1D:   42,
				RS2D:   43,
				Imm:    0b000000001010,
			},
			wantTaken:  true,
			wantOffset: 20,
		},
		{
			name: "blt taken signed neg-pos",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Blt.F3,
				RS1D:   0xFFFFFFF0, // -16
				RS2D:   5,
				Imm:    0b000000001010,
			},
			wantTaken:  true,
			wantOffset: 20,
		},
		{
			name: "bltu not taken signed neg-pos",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Bltu.F3,
				RS1D:   0xFFFFFFF0, // Large unsigned
				RS2D:   5,
				Imm:    0b000000001010,
			},
			wantTaken:  false,
			wantOffset: 4,
		},
		{
			name: "bge taken signed pos-neg",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Bge.F3,
				RS1D:   5,
				RS2D:   0xFFFFFFF0, // -16
				Imm:    0b000000001010,
			},
			wantTaken:  true,
			wantOffset: 20,
		},
		{
			name: "bgeu taken unsigned neg-pos",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Bgeu.F3,
				RS1D:   0xFFFFFFF0, // Large unsigned
				RS2D:   5,
				Imm:    0b000000001010,
			},
			wantTaken:  true,
			wantOffset: 20,
		},
		{
			name: "sign extend offset",
			op: datapath.BranchOp{
				PC:     0x1000,
				Opcode: 0x63,
				Funct3: branch.Beq.F3,
				RS1D:   10,
				RS2D:   10,
				Imm:    0b111111110110, // negative offset
			},
			wantTaken:  true,
			wantOffset: -20,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inCh := make(chan []byte)
			br := New()
			br.Connect(Labels.InBranchOp, inCh)

			if err := br.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				inCh <- datapath.EncodeBranchOp(test.op)
				close(inCh)
			}()

			waiter := make(chan struct{})
			outCh := br.GetPin(Labels.OutOperation)
			go func() {
				defer close(waiter)
				for output := range outCh {
					op := datapath.DecodeOp(output)
					if op.Opcode != test.op.Opcode {
						t.Errorf("expected opcode %x, got %x", test.op.Opcode, op.Opcode)
					}
					if op.AluOp != alu.Ops.Branch1 {
						t.Errorf("expected aluOp %d, got %d", alu.Ops.Branch1, op.AluOp)
					}
					if op.AluOperand1 != test.op.PC {
						t.Errorf("expected operand1 %d, got %d", test.op.PC, op.AluOperand1)
					}
					if int32(op.AluOperand2) != test.wantOffset {
						t.Errorf("expected offset %d, got %d", test.wantOffset, int32(op.AluOperand2))
					}
				}
			}()

			select {
			case <-waiter:
			case <-time.After(50 * time.Millisecond):
				t.Fatal("Brancher took too long")
			}
		})
	}
}

func FuzzBrancher(f *testing.F) {
	f.Add(uint32(10), uint32(10), uint32(0x1000), uint32(0b000000001010), uint8(branch.Beq.F3))
	f.Add(uint32(10), uint32(20), uint32(0x1000), uint32(0b000000001010), uint8(branch.Bne.F3))
	f.Fuzz(func(t *testing.T, rs1d, rs2d, pc, imm uint32, f3 uint8) {
		f3 = f3 % 8
		if f3 != branch.Beq.F3 && f3 != branch.Bne.F3 && f3 != branch.Blt.F3 && f3 != branch.Bge.F3 && f3 != branch.Bltu.F3 && f3 != branch.Bgeu.F3 {
			return
		}

		op := datapath.BranchOp{
			PC:     datapath.XWord(pc),
			Opcode: 0x63,
			Funct3: f3,
			RS1D:   datapath.XWord(rs1d),
			RS2D:   datapath.XWord(rs2d),
			Imm:    imm & 0xFFF, // 12 bits
		}

		offset := (imm & 0xFFF) << 1
		if (offset & 0x1000) != 0 {
			offset |= 0xffffe000
		}

		inCh := make(chan []byte)
		br := New()
		br.Connect(Labels.InBranchOp, inCh)

		if err := br.Run(); err != nil {
			t.Fatal(err)
		}

		go func() {
			inCh <- datapath.EncodeBranchOp(op)
			close(inCh)
		}()

		var taken bool
		switch f3 {
		case branch.Beq.F3:
			taken = rs1d == rs2d
		case branch.Bne.F3:
			taken = rs1d != rs2d
		case branch.Blt.F3:
			taken = int32(rs1d) < int32(rs2d)
		case branch.Bge.F3:
			taken = int32(rs1d) >= int32(rs2d)
		case branch.Bltu.F3:
			taken = rs1d < rs2d
		case branch.Bgeu.F3:
			taken = rs1d >= rs2d
		}

		expectedOffset := int32(4)
		if taken {
			expectedOffset = int32(offset)
		}

		output := <-br.GetPin(Labels.OutOperation)
		res := datapath.DecodeOp(output)

		if int32(res.AluOperand2) != expectedOffset {
			t.Errorf("expected offset %d, got %d", expectedOffset, int32(res.AluOperand2))
		}
	})
}
