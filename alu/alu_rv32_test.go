package alu

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestALUOperations_ToRegister(t *testing.T) {
	tests := []struct {
		name      string
		operation []byte
		regStor   datapath.RegisterData
	}{
		{
			name:      "add,addi: pos-pos",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 7, AluOperand2: 12}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 19},
		},
		{
			name:      "add,addi: pos-neg",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 7, AluOperand2: 0b11111111_11111111_11111111_11111011}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 2},
		},
		{
			// add,addi: -3 + -7
			name:      "add,addi: neg-neg",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 0b11111111111111111111111111111101, AluOperand2: 0b11111111111111111111111111111001}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0b11111111111111111111111111110110},
		},
		{
			name:      "sub: pos-pos",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sub, AluOperand1: 7, AluOperand2: 3}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 4},
		},

		{
			// sub: 7 - (-3)
			name:      "sub: pos-neg",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sub, AluOperand1: 7, AluOperand2: 0b11111111111111111111111111111101}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 10},
		},
		{
			// sub: -7 - (-3)
			name:      "sub: neg-neg",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sub, AluOperand1: 0b11111111111111111111111111111001, AluOperand2: 0b11111111111111111111111111111101}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0b11111111111111111111111111111100},
		},
		{
			name:      "sll",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sll, AluOperand1: 2, AluOperand2: 2}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 8},
		},
		{
			name:      "slli",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sll, AluOperand1: 2, AluOperand2: 3}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 16},
		},
		{
			name:      "slt-true",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Slt, AluOperand1: 2, AluOperand2: 3}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 1},
		},
		{
			// slt 2 < 3
			name:      "slt-true",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Slt, AluOperand1: 2, AluOperand2: 0b11111111111111111111111111111101}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0},
		},
		{
			name:      "slti-false",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Slt, AluOperand1: 5, AluOperand2: 3}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0},
		},
		{
			name:      "sltu-true",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sltu, AluOperand1: 2, AluOperand2: 3}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 1},
		},
		{
			name:      "sltui-false",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sltu, AluOperand1: 5, AluOperand2: 3}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0},
		},
		{
			name:      "xor",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Xor, AluOperand1: 0xFF, AluOperand2: 0b00101}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0xFF ^ 5},
		},
		{
			name:      "xori",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Xor, AluOperand1: 0xFF, AluOperand2: 0b00101}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0xFF ^ 5},
		},
		{
			name:      "and",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.And, AluOperand1: 0xFF, AluOperand2: 0b00101}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0xFF & 5},
		},
		{
			name:      "andi",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.And, AluOperand1: 0xFF, AluOperand2: 0b00101}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0xFF & 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opsCh := make(chan []byte)

			alu := New()
			alu.Connect(Labels.InOperations, opsCh)

			if err := alu.Run(); err != nil {
				t.Fatal(err)
			}

			go func() {
				opsCh <- test.operation
				close(opsCh)
			}()

			waiter := make(chan struct{})
			regDataCh := alu.GetPin(Labels.OutRegData)
			go func() {
				defer close(waiter)
				for output := range regDataCh {
					regData := datapath.DecodeRegStore(output)
					if regData.Value != test.regStor.Value {
						t.Errorf("unexpected regData value: %d", regData.Value)
					}
					if regData.Rd != test.regStor.Rd {
						t.Errorf("unexpected ALU regStor RD: %d", regData.Rd)
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

func TestALU_Run_ToRegister(t *testing.T) {
	opsCh := make(chan []byte)
	alu := New()
	alu.Connect(Labels.InOperations, opsCh)

	// load params
	go func() {
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 7, AluOperand2: 12})
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sll, AluOperand1: 2, AluOperand2: 2})
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Xor, AluOperand1: 0xFF, AluOperand2: 0b00101})
		close(opsCh)
	}()

	if err := alu.Run(); err != nil {
		t.Fatal(err)
	}

	// add
	output := <-alu.GetPin(Labels.OutRegData)
	result := datapath.DecodeRegStore(output)
	if result.Value != 19 {
		t.Errorf("unexpected value %d", result.Value)
	}

	// sll
	output = <-alu.GetPin(Labels.OutRegData)
	result = datapath.DecodeRegStore(output)
	if result.Value != 8 {
		t.Errorf("unexpected value %d", result.Value)
	}

	output = <-alu.GetPin(Labels.OutRegData)
	result = datapath.DecodeRegStore(output)
	if result.Value != 0xFF^5 {
		t.Errorf("unexpected value %d", result.Value)
	}
}
