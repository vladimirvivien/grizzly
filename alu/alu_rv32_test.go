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
		param  []byte
		result datapath.Result
	}{
		{
			name:   "add,addi: pos-pos",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: 0, Funct7: 0, Op1: 7, Op2: 12}),
			result: datapath.Result{Funct3: 0, Rd: 0b00101, AluOut: 19},
		},
		{
			name:   "add,addi: pos-neg",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Add.F3, Funct7: integer.Add.F7, Op1: 7, Op2: 0b11111111_11111111_11111111_11111011}),
			result: datapath.Result{Funct3: integer.Add.F3, Rd: 0b00101, AluOut: 2},
		},
		{
			// add,addi: -3 + -7
			name:   "add,addi: neg-neg",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Addi.F3, Funct7: integer.Addi.F7, Op1: 0b11111111111111111111111111111101, Op2: 0b11111111111111111111111111111001}),
			result: datapath.Result{Funct3: integer.Addi.F3, Rd: 0b00101, AluOut: 0b11111111111111111111111111110110},
		},
		{
			name:   "sub: pos-pos",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sub.F3, Funct7: integer.Sub.F7, Op1: 7, Op2: 3}),
			result: datapath.Result{Funct3: integer.Sub.F3, Rd: 0b00101, AluOut: 4},
		},

		{
			// sub: 7 - (-3)
			name:   "sub: pos-neg",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sub.F3, Funct7: integer.Sub.F7, Op1: 7, Op2: 0b11111111111111111111111111111101}),
			result: datapath.Result{Funct3: integer.Sub.F3, Rd: 0b00101, AluOut: 10},
		},
		{
			// sub: -7 - (-3)
			name:   "sub: neg-neg",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sub.F3, Funct7: integer.Sub.F7, Op1: 0b11111111111111111111111111111001, Op2: 0b11111111111111111111111111111101}),
			result: datapath.Result{Funct3: integer.Sub.F3, Rd: 0b00101, AluOut: 0b11111111111111111111111111111100},
		},
		{
			name:   "sll",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sll.F3, Funct7: integer.Sll.F7, Op1: 2, Op2: 2}),
			result: datapath.Result{Funct3: integer.Sll.F3, Rd: 0b00101, AluOut: 8},
		},
		{
			name:   "slli",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slli.F3, Funct7: integer.Slli.F7, Op1: 2, Op2: 3}),
			result: datapath.Result{Funct3: integer.Slli.F3, Rd: 0b00101, AluOut: 16},
		},
		{
			name:   "slt-true",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slt.F3, Funct7: integer.Slt.F7, Op1: 2, Op2: 3}),
			result: datapath.Result{Funct3: integer.Slt.F3, Rd: 0b00101, AluOut: 1},
		},
		{
			// slt 2 < 3
			name:   "slt-true",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slt.F3, Funct7: integer.Slt.F7, Op1: 2, Op2: 0b11111111111111111111111111111101}),
			result: datapath.Result{Funct3: integer.Slt.F3, Rd: 0b00101, AluOut: 0},
		},
		{
			name:   "slti-false",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Slti.F3, Funct7: integer.Slti.F7, Op1: 5, Op2: 3}),
			result: datapath.Result{Funct3: integer.Slti.F3, Rd: 0b00101, AluOut: 0},
		},
		{
			name:   "sltu-true",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sltu.F3, Funct7: integer.Sltu.F7, Op1: 2, Op2: 3}),
			result: datapath.Result{Funct3: integer.Sltu.F3, Rd: 0b00101, AluOut: 1},
		},
		{
			name:   "sltui-false",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sltiu.F3, Funct7: integer.Sltiu.F7, Op1: 5, Op2: 3}),
			result: datapath.Result{Funct3: integer.Sltiu.F3, Rd: 0b00101, AluOut: 0},
		},
		{
			name:   "xor",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Xor.F3, Funct7: integer.Xor.F7, Op1: 0xFF, Op2: 0b00101}),
			result: datapath.Result{Funct3: integer.Xor.F3, Rd: 0b00101, AluOut: 0xFF ^ 5},
		},
		{
			name:   "xori",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Xori.F3, Funct7: integer.Xori.F7, Op1: 0xFF, Op2: 0b00101}),
			result: datapath.Result{Funct3: integer.Xori.F3, Rd: 0b00101, AluOut: 0xFF ^ 5},
		},
		{
			name:   "and",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.And.F3, Funct7: integer.And.F7, Op1: 0xFF, Op2: 0b00101}),
			result: datapath.Result{Funct3: integer.And.F3, Rd: 0b00101, AluOut: 0xFF & 5},
		},
		{
			name:   "andi",
			param:  datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Andi.F3, Funct7: integer.Andi.F7, Op1: 0xFF, Op2: 0b00101}),
			result: datapath.Result{Funct3: integer.Andi.F3, Rd: 0b00101, AluOut: 0xFF & 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paramsCh := make(chan []byte)

			alu := New()
			alu.Connect(Labels.InParams, paramsCh)

			if err := alu.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
			    paramsCh <- test.param
			    close(paramsCh)
			}()

			waiter := make(chan struct{})
			resultCh := alu.GetPin(Labels.OutResult)
			go func() {
				defer close(waiter)
				for output := range resultCh {
					result := datapath.DecodeResult(output)
					if result.AluOut != test.result.AluOut {
						t.Errorf("unexpected ALU result value: %d", result.AluOut)
					}
					if result.Rd != test.result.Rd {
						t.Errorf("unexpected ALU result RD: %d", result.Rd)
					}
					if result.Funct3 != test.result.Funct3 {
						t.Errorf("unexpected ALU result Funct3: %d", result.Funct3)
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

func TestALU_Run_Multiple(t *testing.T) {
	paramsCh := make(chan []byte)
	alu := New()
	alu.Connect(Labels.InParams, paramsCh)

	// load params
	go func(){
		paramsCh <- datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: 0, Funct7: 0, Op1: 7, Op2: 12})
		paramsCh <- datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Sll.F3, Funct7: integer.Sll.F7, Op1: 2, Op2: 2})
		paramsCh <- datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, Funct3: integer.Xor.F3, Funct7: integer.Xor.F7, Op1: 0xFF, Op2: 0b00101})
		close(paramsCh)
	}()

	if err := alu.Run(); err != nil {
		t.Fatal(err)
	}

	// add
	output := <- alu.GetPin(Labels.OutResult)
	result := datapath.DecodeResult(output)
	if result.AluOut != 19{
		t.Errorf("unexpected value %d", result.AluOut)
	}

	// sll
	output = <- alu.GetPin(Labels.OutResult)
	result = datapath.DecodeResult(output)
	if result.AluOut != 8 {
		t.Errorf("unexpected value %d", result.AluOut)
	}

	output = <- alu.GetPin(Labels.OutResult)
	result = datapath.DecodeResult(output)
	if result.AluOut != 0xFF ^ 5 {
		t.Errorf("unexpected value %d", result.AluOut)
	}
}
