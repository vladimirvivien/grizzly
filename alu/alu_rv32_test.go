//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package alu

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
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
			pcOpCh := alu.GetPin(Labels.OutPcOp)
			go func() {
				defer close(waiter)
				for output := range regDataCh {
					regData := datapath.DecodeRegData(output)
					if regData.Value != test.regStor.Value {
						t.Errorf("unexpected regData value: %d", regData.Value)
					}
					if regData.Rd != test.regStor.Rd {
						t.Errorf("unexpected ALU regStor RD: %d", regData.Rd)
					}
					<-pcOpCh // drain
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

func TestALU_Run_ToMem(t *testing.T) {
	opsCh := make(chan []byte)
	alu := New()
	alu.Connect(Labels.InOperations, opsCh)
	pcCh := alu.GetPin(Labels.OutPcOp)

	// load params
	go func() {
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: isa.Opcodes.L, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 7, AluOperand2: 12})
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: isa.Opcodes.S, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 2, AluOperand2: 2, MemData: 12345})
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: isa.Opcodes.L, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 0xFF, AluOperand2: 0b00101})
		close(opsCh)
	}()

	if err := alu.Run(); err != nil {
		t.Fatal(err)
	}

	// load
	output := <-alu.GetPin(Labels.OutMemOp)
	result := datapath.DecodeMemOp(output)
	if result.Addr != 19 {
		t.Errorf("unexpected value %d", result.Addr)
	}
	<-pcCh // drain pc op


	// store
	output = <-alu.GetPin(Labels.OutMemOp)
	result = datapath.DecodeMemOp(output)
	if result.Addr != 4 {
		t.Errorf("unexpected value %d", result.Addr)
	}
	if result.Data != 12345 {
		t.Errorf("unexpected value %d", result.Data)
	}
	<-pcCh // drain pc op

	// load
	output = <-alu.GetPin(Labels.OutMemOp)
	result = datapath.DecodeMemOp(output)
	if result.Addr != 0xFF+5 {
		t.Errorf("unexpected value %d", result.Addr)
	}
	<-pcCh // drain pc op
}
