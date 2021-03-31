package alu

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	Labels = struct {
		InOperations datapath.Pin
		OutRegData   datapath.Pin
		OutMemOp     datapath.Pin
	}{
		InOperations: datapath.Pin("alu.in.operations"),
		OutRegData:   datapath.Pin("alu.out.reg_data"),
		OutMemOp:     datapath.Pin("alu.out.mem_op"),
	}
)

type ALU struct {
	*datapath.BaseComponent
	outReg chan []byte
	outMem chan []byte
}

func New() *ALU {
	alu := &ALU{
		BaseComponent: datapath.NewBase(),
		outReg:        make(chan []byte),
		outMem:        make(chan []byte),
	}
	alu.Connect(Labels.OutRegData, alu.outReg)
	alu.Connect(Labels.OutMemOp, alu.outMem)
	return alu
}

func (a *ALU) Run() error {
	input := a.GetPin(Labels.InOperations)
	if input == nil {
		return fmt.Errorf("alu: missing input: %s", Labels.InOperations)
	}

	go func() {
		defer close(a.outReg)
		defer close(a.outMem)

		for {
			stream, opened := <-input
			if !opened {
				return
			}

			operation := datapath.DecodeOp(stream)

			var result datapath.XWord
			switch operation.AluOp {
			case
				// addition: add, addi
				Ops.Add:
				result = operation.AluOperand1 + operation.AluOperand2

			case
				// subtraction: sub
				Ops.Sub:
				result = operation.AluOperand1 - operation.AluOperand2

			case
				// shift logical left: sll, slli
				Ops.Sll:
				result = operation.AluOperand1 << operation.AluOperand2

			case
				// set if less then (signed): slt, slti
				Ops.Slt:
				if datapath.SXWord(operation.AluOperand1) < datapath.SXWord(operation.AluOperand2) {
					result = 1
				}

			case
				// set if less then (unsigned): sltu, sltiu
				Ops.Sltu:
				if operation.AluOperand1 < operation.AluOperand2 {
					result = 1
				}

			case
				// xor, xori
				Ops.Xor:
				result = operation.AluOperand1 ^ operation.AluOperand2

			case
				// shift right logical: srl, srli
				Ops.Srl:
				result = operation.AluOperand1 >> operation.AluOperand2

			case
				// shift right arithmetic: sra, srai
				Ops.Sra:
				result = datapath.XWord(datapath.SXWord(operation.AluOperand1) >> operation.AluOperand2)

			case
				// or, ori
				Ops.Or:
				result = operation.AluOperand1 | operation.AluOperand2

			case
				// and, andi
				Ops.And:
				result = operation.AluOperand1 & operation.AluOperand2
			}

			// route alu regStor to other components
			switch operation.Opcode {
			case isa.Opcodes.R, isa.Opcodes.RI:
				a.outReg <- datapath.EncodeRegStore(datapath.RegisterData{Rd: operation.Rd, Value: result})
			case isa.Opcodes.L, isa.Opcodes.S:
				a.outMem <- datapath.EncodeMemOp(datapath.MemOp{
					Opcode: operation.Opcode,
					Rd:     operation.Rd,
					Op:     operation.MemOp,
					Addr:   result,
					Data:   operation.MemData,
				})
			}
		}
	}()

	return nil
}

// mulh** returns high 32-bit portion of multiplication product.
// For 32-bit operands, operation assumes 64-bit host machine.

// mulh interpret operands as signed
func mulh(data1, data2 datapath.XWord) datapath.XWord {
	result := (int64(data1) * int64(data2)) >> datapath.XWordLen
	return datapath.XWord(result)
}

// mulhsu interpret operands as signed/unsigned
func mulhsu(data1, data2 datapath.XWord) datapath.XWord {
	result := (uint64(int32(data1)) * uint64(data2)) >> datapath.XWordLen
	return datapath.XWord(result)
}

func mulhu(data1, data2 datapath.XWord) datapath.XWord {
	result := (uint64(data1) * uint64(data2)) >> datapath.XWordLen
	return datapath.XWord(result)
}
