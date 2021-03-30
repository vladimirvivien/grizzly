package alu

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

var (
	Labels = struct {
		InParams  datapath.Pin
		OutResult datapath.Pin
	}{
		InParams:  datapath.Pin("alu.in.params"),
		OutResult: datapath.Pin("alu.out.result"),
	}
)

type ALU struct {
	*datapath.BaseComponent
	output chan []byte
}

func New() *ALU {
	alu := &ALU{
		BaseComponent: datapath.NewBase(),
		output:        make(chan []byte),
	}
	alu.Connect(Labels.OutResult, alu.output)
	return alu
}

func (a *ALU) Run() error {
	input := a.GetPin(Labels.InParams)
	if input == nil {
		return fmt.Errorf("alu: missing input: %s", Labels.InParams)
	}

	go func() {
		defer close(a.output)
		result := a.output
		for {
			stream, opened := <-input
			if !opened {
				return
			}

			params := datapath.DecodeOp(stream)

			var value datapath.XWord
			switch {
			case
				// addition: add, addi
				params.Funct7 == integer.Add.F7 && params.Funct3 == integer.Add.F3,
				params.Funct7 == integer.Addi.F7 && params.Funct3 == integer.Addi.F3:
				value = params.Op1 + params.Op2

			case
				// subtraction: sub
				params.Funct7 == integer.Sub.F7 && params.Funct3 == integer.Sub.F3:
				value = params.Op1 - params.Op2

			case
				// shift logical left: sll, slli
				params.Funct7 == integer.Sll.F7 && params.Funct3 == integer.Sll.F3,
				params.Funct7 == integer.Slli.F7 && params.Funct3 == integer.Slli.F3:
				value = params.Op1 << params.Op2

			case
				// set if less then (signed): slt, slti
				params.Funct7 == integer.Slt.F7 && params.Funct3 == integer.Slt.F3,
				params.Funct7 == integer.Slti.F7 && params.Funct3 == integer.Slti.F3:
				if datapath.SXWord(params.Op1) < datapath.SXWord(params.Op2) {
					value = 1
				}

			case
				// set if less then (unsigned): sltu, sltiu
				params.Funct7 == integer.Sltu.F7 && params.Funct3 == integer.Sltu.F3,
				params.Funct7 == integer.Sltiu.F7 && params.Funct3 == integer.Sltiu.F3:
				if params.Op1 < params.Op2 {
					value = 1
				}

			case
				// xor, xori
				params.Funct7 == integer.Xor.F7 && params.Funct3 == integer.Xor.F3,
				params.Funct7 == integer.Xori.F7 && params.Funct3 == integer.Xori.F3:
				value = params.Op1 ^ params.Op2

			case
				// shift right logical: srl, srli
				params.Funct7 == integer.Srl.F7 && params.Funct3 == integer.Srl.F3,
				params.Funct7 == integer.Srli.F7 && params.Funct3 == integer.Srli.F3:
				value = params.Op1 >> params.Op2

			case
				// shift right arithmetic: sra, srai
				params.Funct7 == integer.Sra.F7 && params.Funct3 == integer.Sra.F3,
				params.Funct7 == integer.Srai.F7 && params.Funct3 == integer.Srai.F3:
				value = datapath.XWord(datapath.SXWord(params.Op1) >> params.Op2)

			case
				// or, ori
				params.Funct7 == integer.Or.F7 && params.Funct3 == integer.Or.F3,
				params.Funct7 == integer.Ori.F7 && params.Funct3 == integer.Ori.F3:
				value = params.Op1 | params.Op2

			case
				// and, andi
				params.Funct7 == integer.And.F7 && params.Funct3 == integer.And.F3,
				params.Funct7 == integer.Andi.F7 && params.Funct3 == integer.Andi.F3:
				value = params.Op1 & params.Op2
			}

			result <- datapath.EncodeResult(datapath.Result{
				Opcode: params.Opcode,
				Funct3: params.Funct3,
				Funct7: params.Funct7,
				AluOut: value,
				Rd:     params.Rd,
				Data:   params.Data,
			})
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
