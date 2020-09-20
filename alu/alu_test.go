package alu

import (
	"testing"
	"time"
)

func TestALUROps(t *testing.T) {
	d1wires := make(chan uint32)
	d2wires := make(chan uint32)
	opwires := make(chan uint32)

	tests := []struct {
		name       string
		operands   func() (uint32, uint32)
		aluOp      uint32
		expected   func() uint32
		shouldFail bool
	}{
		{
			name:     "add-pos-pos",
			operands: func() (uint32, uint32) { return 0x2, 0x5 },
			expected: func() uint32 { return 0x7 },
			aluOp:    Ops.Add,
		},
		{
			name:     "add-pos-neg",
			operands: func() (uint32, uint32) { return 0x2, 0b11111111_11111111_11111111_11111001 },
			expected: func() uint32 { return uint32(0b11111111_11111111_11111111_11111011) },
			aluOp:    Ops.Add,
		},
		{
			name: "add-neg-neg",
			operands: func() (uint32, uint32) {
				return 0b11111111_11111111_11111111_11111001, 0b11111111_11111111_11111111_11111011
			},
			expected: func() uint32 { return uint32(0b11111111_11111111_11111111_11110100) },
			aluOp:    Ops.Add,
		},
		{
			name:     "sub",
			operands: func() (uint32, uint32) { return 0x7, 0x3 },
			expected: func() uint32 { return 0x4 },
			aluOp:    Ops.Sub,
		},
		{
			name:     "sll",
			operands: func() (uint32, uint32) { return 0x2, 0x2 },
			expected: func() uint32 { return 0x8 },
			aluOp:    Ops.Sll,
		},
		{
			name:     "slt-true", // true
			operands: func() (uint32, uint32) { return 0x4, 0x1C },
			expected: func() uint32 { return 0x1 },
			aluOp:    Ops.Slt,
		},
		{
			name:     "slt-false", // false
			operands: func() (uint32, uint32) { return 0x1C, 0x10 },
			expected: func() uint32 { return 0x0 },
			aluOp:    Ops.Slt,
		},
		{
			name:     "sltu-true",
			operands: func() (uint32, uint32) { return 0x1C, 0xFF },
			expected: func() uint32 { return 0x1 },
			aluOp:    Ops.Sltu,
		},
		{
			name:     "sltu-false",
			operands: func() (uint32, uint32) { return 0x1C, 0x10 },
			expected: func() uint32 { return 0x0 },
			aluOp:    Ops.Sltu,
		},

		// ***** Multiplication Tests *******
		// mul
		{
			name:     "mul-pos-pos",
			operands: func() (uint32, uint32) { return 0x2, 0x5 },
			expected: func() uint32 { return 0xA },
			aluOp:    Ops.Mul,
		},
		{
			name:     "mul-pos-neg",
			operands: func() (uint32, uint32) { return 0x2, 0b11111111_11111111_11111111_11111011 },
			expected: func() uint32 { return 0b11111111_11111111_11111111_11110110 },
			aluOp:    Ops.Mul,
		},
		{
			name: "mul-neg-neg",
			operands: func() (uint32, uint32) {
				return 0b11111111_11111111_11111111_11111110, 0b11111111_11111111_11111111_11111011
			},
			expected: func() uint32 { return 0xA },
			aluOp:    Ops.Mul,
		},

		// Mulh
		{
			name:     "mulh-pos-pos small",
			operands: func() (uint32, uint32) { return 0x2, 0x5 },
			expected: func() uint32 { return 0 },
			aluOp:    Ops.Mulh,
		},
		{
			name:     "mulh-pos-pos large",
			operands: func() (uint32, uint32) { return 1100200, 1200300 },
			expected: func() uint32 { return 307 },
			aluOp:    Ops.Mulh,
		},
		{
			name:     "mulh-pos-neg large",
			operands: func() (uint32, uint32) { return 1100200, 0b11111111_11101101_10101111_01010100 },
			expected: func() uint32 { return 0b100001100100001110100 },
			aluOp:    Ops.Mulh,
		},
		{
			name: "mulh-neg-neg large",
			operands: func() (uint32, uint32) {
				return 0b11111111_11101111_00110110_01011000, 0b11111111_11101101_10101111_01010100
			},
			expected: func() uint32 { return 0b11111111110111001110011011011111 },
			aluOp:    Ops.Mulh,
		},

		// Mulhsu
		{
			name:     "mulhsu-pos-pos small",
			operands: func() (uint32, uint32) { return 0x2, 0x5 },
			expected: func() uint32 { return 0 },
			aluOp:    Ops.Mulhsu,
		},
		{
			name:     "mulhsu-pos-pos large",
			operands: func() (uint32, uint32) { return 1100200, 1200300 },
			expected: func() uint32 { return 307 },
			aluOp:    Ops.Mulhsu,
		},
		{
			name:     "mulhsu-pos-neg large",
			operands: func() (uint32, uint32) { return 1100200, 0b11111111_11101101_10101111_01010100 },
			expected: func() uint32 { return 0b100001100100001110100 },
			aluOp:    Ops.Mulhsu,
		},
		{
			name: "mulhsu-neg-neg large",
			operands: func() (uint32, uint32) {
				return 0b11111111_11101111_00110110_01011000, 0b11111111_11101101_10101111_01010100
			},
			expected: func() uint32 { return 0b11111111111011110011011110001011 },
			aluOp:    Ops.Mulhsu,
		},

		// Mulhu
		{
			name:     "mulhu-pos-pos small",
			operands: func() (uint32, uint32) { return 0x2, 0x5 },
			expected: func() uint32 { return 0 },
			aluOp:    Ops.Mulhu,
		},
		{
			name:     "mulhu-pos-pos large",
			operands: func() (uint32, uint32) { return 1100200, 1200300 },
			expected: func() uint32 { return 307 },
			aluOp:    Ops.Mulhu,
		},
		{
			name:     "mulhu-pos-neg large",
			operands: func() (uint32, uint32) { return 1100200, 0b11111111_11101101_10101111_01010100 },
			expected: func() uint32 { return 0b100001100100001110100 },
			aluOp:    Ops.Mulhu,
		},
		{
			name: "mulhu-neg-neg large",
			operands: func() (uint32, uint32) {
				return 0b11111111_11101111_00110110_01011000, 0b11111111_11101101_10101111_01010100
			},
			expected: func() uint32 { return 0b11111111110111001110011011011111 },
			aluOp:    Ops.Mulhu,
		},
	}

	alu := newAlu()

	// wire ports
	alu.SetPin(In.Operand1, d1wires)
	alu.SetPin(In.Operand2, d2wires)
	alu.SetPin(In.Operation, opwires)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})

			go func() {
				op1, op2 := test.operands()
				d1wires <- op1
				d2wires <- op2
				opwires <- test.aluOp
			}()

			if err := alu.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				defer close(wait)
				// interpret as signed for proper comparison
				result := int32(<-alu.GetPin(Out.Result))
				expected := int32(test.expected())
				if result != expected {
					t.Errorf("unexpected test result for op %s: %d", test.name, result)
				}
			}()

			// detect stuck path
			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatalf("ALU operation %s took too long to comlete", test.name)
			}
		})
	}
}
