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
	// internal wires
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

	// Input Loop
	// This loop processes incoming alu operations
	// and places result on internal channels to be sent as output.
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

			// route alu result routing
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
				// store RD <- PC+4
				a.xfrReg <- datapath.EncodeRegData(datapath.RegisterData{Rd: operation.Rd, Value: operation.PC + 4})
				// set next PC = op1+op2
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC: result})
			case isa.Opcodes.JI:
				// store RD <- PC+4
				a.xfrReg <- datapath.EncodeRegData(datapath.RegisterData{Rd: operation.Rd, Value: operation.PC + 4})
				// send next PC = (op1+op2) & 0xffff_fffe
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC: result & 0xffff_fffe})
			case isa.Opcodes.B:
				a.xfrPc <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC: result})
			}
		}
	}()

	// Reg Op Output Loop
	// Sends out Register operations
	go func() {
		defer close(a.outReg)
		for stream := range a.xfrReg {
			a.outReg <- stream
		}
	}()

	// Mem Op Output Loop
	// Sends out Memory operations
	go func() {
		defer close(a.outMem)
		for stream := range a.xfrMem {
			a.outMem <- stream
		}
	}()

	// PC Op Output Loop
	// Sends out program counter operations
	go func() {
		defer close(a.outPc)
		for stream := range a.xfrPc {
			a.outPc <- stream
		}
	}()
	return nil
}

// aluFunc carries out the ALU operations
func aluFunc(operation datapath.Operation) (result datapath.XWord) {
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

	case Ops.Mul:
		result = operation.AluOperand1 * operation.AluOperand2

	case Ops.Mulh:
		result = mulh(operation.AluOperand1, operation.AluOperand2)

	case Ops.Mulhsu:
		result = mulhsu(operation.AluOperand1, operation.AluOperand2)

	case Ops.Mulhu:
		result = mulhu(operation.AluOperand1, operation.AluOperand2)

	case Ops.Div:
		op1 := int32(operation.AluOperand1)
		op2 := int32(operation.AluOperand2)
		if op2 == 0 {
			result = 0xffffffff
		} else if op1 == -2147483648 && op2 == -1 {
			result = uint32(op1)
		} else {
			result = uint32(op1 / op2)
		}

	case Ops.Divu:
		op1 := operation.AluOperand1
		op2 := operation.AluOperand2
		if op2 == 0 {
			result = 0xffffffff
		} else {
			result = op1 / op2
		}

	case Ops.Rem:
		op1 := int32(operation.AluOperand1)
		op2 := int32(operation.AluOperand2)
		if op2 == 0 {
			result = uint32(op1)
		} else if op1 == -2147483648 && op2 == -1 {
			result = 0
		} else {
			result = uint32(op1 % op2)
		}

	case Ops.Remu:
		op1 := operation.AluOperand1
		op2 := operation.AluOperand2
		if op2 == 0 {
			result = op1
		} else {
			result = op1 % op2
		}

	case Ops.Branch1:
		// if branch, result = PC + SigExt(Imm)
		// if not branch, result = PC+4
		result = operation.AluOperand1 + operation.AluOperand2
	}
	return
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
