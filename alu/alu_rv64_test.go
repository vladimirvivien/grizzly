//go:build rv64 || rv64i

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
			name:      "add: pos-pos 64-bit",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 7, AluOperand2: 12}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 19},
		},
		{
			name:      "add: 64-bit sign extension pos-neg",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 7, AluOperand2: 0xFFFFFFFFFFFFFFFB}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 2},
		},
		{
			name:      "add: 64-bit overflow check",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 0xFFFFFFFFFFFFFFFF, AluOperand2: 1}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0},
		},
		{
			name:      "sub: pos-pos",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sub, AluOperand1: 7, AluOperand2: 3}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 4},
		},
		{
			name:      "sub: 64-bit sign-extended result",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sub, AluOperand1: 3, AluOperand2: 7}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0xFFFFFFFFFFFFFFFC},
		},
		{
			name:      "sll 64-bit",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sll, AluOperand1: 2, AluOperand2: 34}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 2 << 34},
		},
		{
			name:      "slt-true 64-bit",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Slt, AluOperand1: 0xFFFFFFFFFFFFFFFF, AluOperand2: 1}), // -1 < 1
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 1},
		},
		{
			name:      "sltu-false 64-bit",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sltu, AluOperand1: 0xFFFFFFFFFFFFFFFF, AluOperand2: 1}), // MaxUint64 < 1
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0},
		},
		{
			name:      "sra 64-bit",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Sra, AluOperand1: 0x8000000000000000, AluOperand2: 1}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0xC000000000000000},
		},
		{
			name:      "xor 64-bit",
			operation: datapath.EncodeOp(datapath.Operation{Opcode: 0b0110011, Rd: 0b00101, AluOp: Ops.Xor, AluOperand1: 0xFFFFFFFFFFFFFFFF, AluOperand2: 0x5555555555555555}),
			regStor:   datapath.RegisterData{Rd: 0b00101, Value: 0xAAAAAAAAAAAAAAAA},
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
						t.Errorf("unexpected regData value: %x, expected %x", regData.Value, test.regStor.Value)
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

	go func() {
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: isa.Opcodes.L, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 7, AluOperand2: 12})
		opsCh <- datapath.EncodeOp(datapath.Operation{Opcode: isa.Opcodes.S, Rd: 0b00101, AluOp: Ops.Add, AluOperand1: 2, AluOperand2: 2, MemData: 12345})
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
	<-pcCh

	// store
	output = <-alu.GetPin(Labels.OutMemOp)
	result = datapath.DecodeMemOp(output)
	if result.Addr != 4 {
		t.Errorf("unexpected value %d", result.Addr)
	}
	if result.Data != 12345 {
		t.Errorf("unexpected value %d", result.Data)
	}
	<-pcCh
}

func FuzzALU(f *testing.F) {
	f.Add(uint64(10), uint64(5), uint8(Ops.Add))
	f.Add(uint64(10), uint64(5), uint8(Ops.Sub))
	f.Fuzz(func(t *testing.T, op1 uint64, op2 uint64, aluOp uint8) {
		if aluOp >= Ops.Branch1 {
			return
		}
		
		operation := datapath.Operation{
			AluOp:       aluOp,
			AluOperand1: op1,
			AluOperand2: op2,
		}

		// Prevent shift amount overflows in shift operations
		if aluOp == Ops.Sll || aluOp == Ops.Srl || aluOp == Ops.Sra {
			operation.AluOperand2 = op2 % 64
		}

		result := aluFunc(operation)

		var expected uint64
		switch aluOp {
		case Ops.Add:
			expected = op1 + op2
		case Ops.Sub:
			expected = op1 - op2
		case Ops.Sll:
			expected = op1 << (op2 % 64)
		case Ops.Slt:
			if int64(op1) < int64(op2) {
				expected = 1
			} else {
				expected = 0
			}
		case Ops.Sltu:
			if op1 < op2 {
				expected = 1
			} else {
				expected = 0
			}
		case Ops.Xor:
			expected = op1 ^ op2
		case Ops.Srl:
			expected = op1 >> (op2 % 64)
		case Ops.Sra:
			expected = uint64(int64(op1) >> (op2 % 64))
		case Ops.Or:
			expected = op1 | op2
		case Ops.And:
			expected = op1 & op2
		default:
			return
		}

		if result != expected {
			t.Errorf("op=%d: op1=%x, op2=%x: got %x, expected %x", aluOp, op1, op2, result, expected)
		}
	})
}
