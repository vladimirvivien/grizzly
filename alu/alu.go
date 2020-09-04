package alu

import (
	"github.com/vladimirvivien/grizzly/device"
)

var (
	Operations = struct {
		Add  uint32
		And  uint32
		Sub  uint32
		Or   uint32
		Sll  uint32
		Slt  uint32
		Sltu uint32
		Sra  uint32
		Srl  uint32
		Xor  uint32

		// mul
		Mul    uint32
		Mulh   uint32
		Mulhsu uint32
		Mulhu  uint32

		// div
		Div  uint32
		Divu uint32
		Rem  uint32
		Remu uint32
	}{
		Add:  0b00000,
		And:  0b00001,
		Sub:  0b00010,
		Or:   0b00011,
		Sll:  0b00100,
		Slt:  0b00101,
		Sltu: 0b00110,
		Sra:  0b00111,
		Srl:  0b01000,
		Xor:  0b01001,

		// Mul
		Mul:    0b01010,
		Mulh:   0b01011,
		Mulhsu: 0b01100,
		Mulhu:  0b01101,

		// Div
		Div:  0b01110,
		Divu: 0b01111,
		Rem:  0b10000,
		Remu: 0b10001,
	}
)

var (
	In = struct {
		Operand1  device.PinLabel // operand1
		Operand2  device.PinLabel // operand2
		Operation device.PinLabel // function bits
	}{
		Operand1:  "alu.data1.in",
		Operand2:  "alu.data2.in",
		Operation: "alu.operation.in",
	}

	Out = struct {
		Result device.PinLabel
		Zero   device.PinLabel
	}{
		Result: "alu.result.out",
		Zero:   "alu.zero.out",
	}
)

type ALU struct {
	resultOut device.Wires // output
	zeroOut   device.Wires // zero line
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
			// ALU Operation
			op := <-a.GetPin(In.Operation)

			switch op {
			// addition: add, addi
			case Operations.Add:
				a.resultOut <- data1 + data2

			// sub
			case Operations.Sub:
				a.resultOut <- data1 - data2

			// shift logical left: sll, slli
			case Operations.Sll:
				a.resultOut <- data1 << data2

			// set if less then (signed): slt, slti
			case Operations.Slt:
				var result uint32
				if int32(data1) < int32(data2) {
					result = 1
				}
				a.resultOut <- result

			// set if less then (unsigned): sltu, sltiu
			case Operations.Sltu:
				var result uint32
				if data1 < data2 {
					result = 1
				}
				a.resultOut <- result

			// or, ori
			case Operations.Xor:
				a.resultOut <- data1 ^ data2

			// shift right logical: srl, srli
			case Operations.Srl:
				a.resultOut <- data1 >> data2

			// shift right arithmetic: sra, srai
			case Operations.Sra:
				a.resultOut <- uint32(int32(data1) >> data2)

			// or, ori
			case Operations.Or:
				a.resultOut <- data1 | data2

			// and, andi
			case Operations.And:
				a.resultOut <- data1 & data2

			// mul
			case Operations.Mul:
				a.resultOut <- data1 * data2
			case Operations.Mulh:
				a.resultOut <- mulh(data1, data2)
			case Operations.Mulhsu:
				a.resultOut <- mulhsu(data1, data2)
			case Operations.Mulhu:
				a.resultOut <- mulhu(data1, data2)

			case Operations.Div:
			case Operations.Divu:
			case Operations.Rem:
			case Operations.Remu:
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
