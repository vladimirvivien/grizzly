package alu

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/isa"
)

func TestALUROps(t *testing.T) {
	d1wires := make(chan uint32)
	d2wires := make(chan uint32)
	opwires := make(chan uint32)

	tests := []struct {
		name     string
		operands func() (uint32, uint32)
		aluOp    uint32
		expected func() uint32
	}{
		{
			name:     "add-pos-pos",
			operands: func() (uint32, uint32) { return 0x2, 0x5 },
			expected: func() uint32 { return 0x7 },
			aluOp:    isa.Add.Functs,
		},
		{
			name:     "add-pos-neg",
			operands: func() (uint32, uint32) { return 0x2, 0b11111111_11111111_11111111_11111001 },
			expected: func() uint32 { return uint32(0b11111111_11111111_11111111_11111011) },
			aluOp:    isa.Add.Functs,
		},
		{
			name: "add-neg-neg",
			operands: func() (uint32, uint32) {
				return 0b11111111_11111111_11111111_11111001, 0b11111111_11111111_11111111_11111011
			},
			expected: func() uint32 { return uint32(0b11111111_11111111_11111111_11110100) },
			aluOp:    isa.Add.Functs,
		},
		{
			name:     "sub",
			operands: func() (uint32, uint32) { return 0x7, 0x3 },
			expected: func() uint32 { return 0x4 },
			aluOp:    isa.Sub.Functs,
		},
		{
			name:     "sll",
			operands: func() (uint32, uint32) { return 0x2, 0x2 },
			expected: func() uint32 { return 0x8 },
			aluOp:    isa.Sll.Functs,
		},
		{
			name:     "slt-true", // true
			operands: func() (uint32, uint32) { return 0x4, 0x1C },
			expected: func() uint32 { return 0x1 },
			aluOp:    isa.Slt.Functs,
		},
		{
			name:     "slt-false", // false
			operands: func() (uint32, uint32) { return 0x1C, 0x10 },
			expected: func() uint32 { return 0x0 },
			aluOp:    isa.Slt.Functs,
		},
		{
			name:     "sltu-true",
			operands: func() (uint32, uint32) { return 0x1C, 0xFF },
			expected: func() uint32 { return 0x1 },
			aluOp:    isa.Sltu.Functs,
		},
		{
			name:     "sltu-false",
			operands: func() (uint32, uint32) { return 0x1C, 0x10 },
			expected: func() uint32 { return 0x0 },
			aluOp:    isa.Sltu.Functs,
		},
	}

	alu := newAlu()

	// wire ports
	alu.Data1In(d1wires)
	alu.Data2In(d2wires)
	alu.FunctsIn(opwires)

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
				result := int32(<-alu.DataOut())
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
