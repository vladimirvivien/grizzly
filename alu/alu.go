package alu

import (
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	In = struct {
		Operand1 device.PinLabel // operand1
		Operand2 device.PinLabel // operand2
		Functs   device.PinLabel // function bits
	}{
		Operand1: "alu.data1.in",
		Operand2: "alu.data2.in",
		Functs:   "alu.funct.in",
	}

	Out = struct {
		Result device.PinLabel
	}{
		Result: "alu.result.out",
	}
)

type ALU struct {
	resultOut device.Wires // output
	*device.Base
}

func New() device.Type {
	return newAlu()
}

func newAlu() *ALU {
	a := &ALU{
		resultOut: device.MakeWires(),
		Base:      device.NewBase(),
	}

	a.SetPin(Out.Result, a.resultOut)

	return a
}

// Run starts the ALU.
// Data1 and Data2 are read sequentially and must
// be available or risk blocking.
func (a *ALU) Run() error {
	go func() {
		defer close(a.resultOut)
		for {
			data1 := <-a.GetPin(In.Operand1)
			data2 := <-a.GetPin(In.Operand2)
			// functs concatanate all function bits
			functs := <-a.GetPin(In.Functs)

			switch functs {
			// addition: add, addi
			case isa.Add.Functs, isa.Addi.Functs:
				a.resultOut <- data1 + data2

			// sub
			case isa.Sub.Functs:
				a.resultOut <- data1 - data2

			// shift logical left: ssl, ssli
			case isa.Sll.Functs, isa.Slli.Functs:
				a.resultOut <- data1 << data2

			// set if less then (signed): slt, slti
			case isa.Slt.Functs, isa.Slti.Functs:
				var result uint32
				if int32(data1) < int32(data2) {
					result = 1
				}
				a.resultOut <- result

			// set if less then (unsigned): sltu, sltiu
			case isa.Sltu.Functs, isa.Sltiu.Functs:
				var result uint32
				if data1 < data2 {
					result = 1
				}
				a.resultOut <- result

			// or, ori
			case isa.Or.Functs, isa.Ori.Functs:
				a.resultOut <- data1 ^ data2

			// shift right logical: srl, srli
			case isa.Srl.Functs, isa.Srli.Functs:
				a.resultOut <- data1 >> data2

			// shift right arithmetic: sra, srai
			case isa.Sra.Functs, isa.Srai.Functs:
				a.resultOut <- uint32(int32(data1) >> data2)

			// or, ori
			case isa.Or.Functs, isa.Ori.Functs:
				a.resultOut <- data1 | data2

			// and, andi
			case isa.And.Functs, isa.Andi.Functs:
				a.resultOut <- data1 & data2

			// mul
			case isa.Mul.Functs:
				a.resultOut <- data1 * data2
			case isa.Mulh.Functs:
				a.resultOut <- mulh(data1, data2)
			case isa.Mulhsu.Functs:
				a.resultOut <- mulhsu(data1, data2)
			case isa.Mulhu.Functs:
				a.resultOut <- mulhu(data1, data2)

			case isa.Div.Functs:
			case isa.Divu.Functs:
			case isa.Rem.Functs:
			case isa.Remu.Functs:
			}

		}
	}()

	return nil
}

// mulh** returns high 32-bit portion of multiplication product.
// For 32-bit operands, operation assumes 64-bit host machine.

// mulh interpret operands as signed
func mulh(data1, data2 uint32) uint32 {
	result := (int64(data1) * int64(data2)) >> 32
	return uint32(result)
}

// mulhsu interpret operands as signed/unsigned
func mulhsu(data1, data2 uint32) uint32 {
	result := (uint64(int32(data1)) * uint64(data2)) >> 32
	return uint32(result)
}

func mulhu(data1, data2 uint32) uint32 {
	result := (uint64(data1) * uint64(data2)) >> 32
	return uint32(result)
}
