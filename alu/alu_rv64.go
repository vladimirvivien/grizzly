//go:build rv64 || rv64i

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
		OutPcOp      datapath.Pin
	}{
		InOperations: datapath.Pin("alu.in.operations"),
		OutRegData:   datapath.Pin("alu.out.reg_data"),
		OutMemOp:     datapath.Pin("alu.out.mem_op"),
		OutPcOp:      datapath.Pin("alu.out.pc_op"),
	}
)

type ALU struct {
	*datapath.BaseComponent
	outReg chan []byte
	outMem chan []byte
	outPc  chan []byte
	xfrReg chan []byte
	xfrMem chan []byte
	xfrPc  chan []byte
}

func New() *ALU {
	alu := &ALU{
		BaseComponent: datapath.NewBase(),
		outReg:        make(chan []byte),
		outMem:        make(chan []byte),
		outPc:         make(chan []byte),
		xfrReg:        make(chan []byte),
		xfrMem:        make(chan []byte),
		xfrPc:         make(chan []byte),
	}
	alu.Connect(Labels.OutRegData, alu.outReg)
	alu.Connect(Labels.OutMemOp, alu.outMem)
	alu.Connect(Labels.OutPcOp, alu.outPc)
	return alu
}

func (a *ALU) Run() error {
	input := a.GetPin(Labels.InOperations)
	if input == nil {
		return fmt.Errorf("alu: missing input: %s", Labels.InOperations)
	}

	go func() {
		defer close(a.xfrReg)
		defer close(a.xfrMem)
		defer close(a.xfrPc)

		for {
			stream, opened := <-input
			if !opened {
				return
			}

			operation := datapath.DecodeOp(stream)
			result := aluFunc(operation)

			switch operation.Opcode {
			case isa.Opcodes.R, isa.Opcodes.RI:
				a.xfrReg <- datapath.EncodeRegData(datapath.RegisterData{Rd: operation.Rd, Value: result})
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 0, PC: 0})
			case isa.Opcodes.L, isa.Opcodes.S:
				a.xfrMem <- datapath.EncodeMemOp(datapath.MemOp{
					Opcode: operation.Opcode,
					Rd:     operation.Rd,
					Op:     operation.MemOp,
					Addr:   result,
					Data:   operation.MemData,
				})
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 0, PC: 0})
			case isa.Opcodes.J:
				a.xfrReg <- datapath.EncodeRegData(datapath.RegisterData{Rd: operation.Rd, Value: operation.PC + 4})
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC: result})
			case isa.Opcodes.JI:
				a.xfrReg <- datapath.EncodeRegData(datapath.RegisterData{Rd: operation.Rd, Value: operation.PC + 4})
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC: result & 0xfffffffffffffffe})
			case isa.Opcodes.B:
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC: result})
			}
		}
	}()

	go func() {
		defer close(a.outReg)
		for stream := range a.xfrReg {
			a.outReg <- stream
		}
	}()

	go func() {
		defer close(a.outMem)
		for stream := range a.xfrMem {
			a.outMem <- stream
		}
	}()

	go func() {
		defer close(a.outPc)
		for stream := range a.xfrPc {
			a.outPc <- stream
		}
	}()
	return nil
}

var extOperations = make(map[uint8]func(op1, op2 datapath.XWord) datapath.XWord)

func registerExtOp(op uint8, fn func(op1, op2 datapath.XWord) datapath.XWord) {
	extOperations[op] = fn
}

func aluFunc(operation datapath.Operation) (result datapath.XWord) {
	switch operation.AluOp {
	case Ops.Add:
		result = operation.AluOperand1 + operation.AluOperand2
	case Ops.Sub:
		result = operation.AluOperand1 - operation.AluOperand2
	case Ops.Sll:
		result = operation.AluOperand1 << operation.AluOperand2
	case Ops.Slt:
		if datapath.SXWord(operation.AluOperand1) < datapath.SXWord(operation.AluOperand2) {
			result = 1
		}
	case Ops.Sltu:
		if operation.AluOperand1 < operation.AluOperand2 {
			result = 1
		}
	case Ops.Xor:
		result = operation.AluOperand1 ^ operation.AluOperand2
	case Ops.Srl:
		result = operation.AluOperand1 >> operation.AluOperand2
	case Ops.Sra:
		result = datapath.XWord(datapath.SXWord(operation.AluOperand1) >> operation.AluOperand2)
	case Ops.Or:
		result = operation.AluOperand1 | operation.AluOperand2
	case Ops.And:
		result = operation.AluOperand1 & operation.AluOperand2
	case Ops.Branch1:
		result = operation.AluOperand1 + operation.AluOperand2
	default:
		if fn, exists := extOperations[operation.AluOp]; exists {
			result = fn(operation.AluOperand1, operation.AluOperand2)
		} else {
			panic(fmt.Sprintf("unknown alu operation: %d", operation.AluOp))
		}
	}
	return
}
