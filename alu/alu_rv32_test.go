package alu

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

func TestALUOperations(t *testing.T) {
	tests := []struct {
		name   string
		param  datapath.AluParam
		result datapath.AluResult
	}{
		{
			name:   "add,addi: pos-pos",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: 0, Funct7: 0, Op1: 7, Op2: 12},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 19},
		},
		{
			name:   "add,addi: pos-neg",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Add.F3, Funct7: integer.Add.F7, Op1: 7, Op2: 0b11111111_11111111_11111111_11111011},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 2},
		},
		{
			// add,addi: -3 + -7
			name:   "add,addi: neg-neg",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Addi.F3, Funct7: integer.Addi.F7, Op1: 0b11111111111111111111111111111101, Op2: 0b11111111111111111111111111111001},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0b11111111111111111111111111110110},
		},
		{
			name:   "sub: pos-pos",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sub.F3, Funct7: integer.Sub.F7, Op1: 7, Op2: 3},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 4},
		},

		{
			// sub: 7 - (-3)
			name:   "sub: pos-neg",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sub.F3, Funct7: integer.Sub.F7, Op1: 7, Op2: 0b11111111111111111111111111111101},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 10},
		},
		{
			// sub: -7 - (-3)
			name:   "sub: neg-neg",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sub.F3, Funct7: integer.Sub.F7, Op1: 0b11111111111111111111111111111001, Op2: 0b11111111111111111111111111111101},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0b11111111111111111111111111111100},
		},
		{
			name:   "sll",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sll.F3, Funct7: integer.Sll.F7, Op1: 2, Op2: 2},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 8},
		},
		{
			name:   "slli",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slli.F3, Funct7: integer.Slli.F7, Op1: 2, Op2: 3},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 16},
		},
		{
			name:   "slt-true",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slt.F3, Funct7: integer.Slt.F7, Op1: 2, Op2: 3},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 1},
		},
		{
			// slt 2 < 3
			name:   "slt-true",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slt.F3, Funct7: integer.Slt.F7, Op1: 2, Op2: 0b11111111111111111111111111111101},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0},
		},
		{
			name:   "slti-false",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slti.F3, Funct7: integer.Slti.F7, Op1: 5, Op2: 3},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0},
		},
		{
			name:   "sltu-true",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sltu.F3, Funct7: integer.Sltu.F7, Op1: 2, Op2: 3},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 1},
		},
		{
			name:   "sltui-false",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sltiu.F3, Funct7: integer.Sltiu.F7, Op1: 5, Op2: 3},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0},
		},
		{
			name:   "xor",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Xor.F3, Funct7: integer.Xor.F7, Op1: 0xFF, Op2: 0b00101},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0xFF ^ 5},
		},
		{
			name:   "xori",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Xori.F3, Funct7: integer.Xori.F7, Op1: 0xFF, Op2: 0b00101},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0xFF ^ 5},
		},
		{
			name:   "and",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.And.F3, Funct7: integer.And.F7, Op1: 0xFF, Op2: 0b00101},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0xFF & 5},
		},
		{
			name:   "andi",
			param:  datapath.AluParam{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Andi.F3, Funct7: integer.Andi.F7, Op1: 0xFF, Op2: 0b00101},
			result: datapath.AluResult{F3: 0, Rd: 0b00101, Value: 0xFF & 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paramsCh := make(chan datapath.AluParam)

			alu := New()
			alu.ParamsInput(paramsCh)

			if err := alu.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
			    paramsCh <- test.param
			    close(paramsCh)
			}()

			waiter := make(chan struct{})
			resultCh := alu.ResultOutput()
			go func() {
				defer close(waiter)
				for result := range resultCh {
					if result.Value != test.result.Value {
						t.Errorf("unexpected ALU result value: %d", result.Value)
					}
					if result.Rd != test.result.Rd {
						t.Errorf("unexpected ALU result RD: %d", result.Rd)
					}
					if result.F3 != test.result.F3 {
						t.Errorf("unexpected ALU result Funct3: %d", result.Rd)
					}
				}
			}()

			select {
			case <-waiter:
			case <-time.After(5 * time.Millisecond):
				t.Fatal("ALU operations took too long to complete")
			}
		})
	}
}
