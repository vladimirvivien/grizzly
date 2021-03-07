package alu

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

type ALU struct {
	paramsInput  <-chan datapath.AluParam
	resultOutput chan datapath.AluResult
}

func New() *ALU {
	return &ALU{
		resultOutput: make(chan datapath.AluResult),
	}
}

func (a *ALU) ParamsInput(ch <-chan datapath.AluParam) {
	a.paramsInput = ch
}

func (a *ALU) ResultOutput() <-chan datapath.AluResult {
	return a.resultOutput
}

func (a *ALU) Run() error {
	if a.paramsInput == nil {
		return fmt.Errorf("alu: missing params input")
	}

	go func() {
		defer close(a.resultOutput)
		result := a.resultOutput
		for {
			params, opened := <-a.paramsInput
			if !opened {
				return
			}

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

			result <- datapath.AluResult{Value: value, Rd: params.Rd}
		}
	}()

	return nil
}
