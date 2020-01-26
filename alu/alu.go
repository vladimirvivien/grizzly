package alu

import (
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/isa"
)

var (
	Wires = struct {
		Data1In,
		Data2In,
		FunctIn,
		DataOut string
	}{
		Data1In: "alu.data1.in",
		Data2In: "alu.data2.in",
		FunctIn: "alu.funct.in",
		DataOut: "alu.data.out",
	}
)

type ALU struct {
	data1In  device.WiresIn // operand 1
	data2In  device.WiresIn // operand 2
	dataOut  device.Wires   // output
	functsIn device.WiresIn // 10-bit inst[31:25]inst[14:12]
}

func New() device.Type {
	return newAlu()
}

func newAlu() *ALU {
	return &ALU{
		dataOut: device.MakeWires(),
	}
}

// Run starts the ALU.
// Data1 and Data2 are read sequentially and must
// be available or risk blocking.
func (a *ALU) Run() error {
	go func() {
		defer close(a.dataOut)
		for {
			data1 := <-a.data1In
			data2 := <-a.data2In
			functs := <-a.functsIn

			switch functs {
			// add
			case isa.Add.Functs:
				a.dataOut <- data1 + data2

			// sub
			case isa.Sub.Functs:
				a.dataOut <- data1 - data2

			// sll - shift logical left
			case isa.Sll.Functs:
				a.dataOut <- data1 << data2

			// slt - set if less then (signed)
			case isa.Slt.Functs:
				var result uint32
				if int32(data1) < int32(data2) {
					result = 1
				}
				a.dataOut <- result

			// sltu - set if less then (unsigned)
			case isa.Sltu.Functs:
				var result uint32
				if data1 < data2 {
					result = 1
				}
				a.dataOut <- result

			// xor
			case isa.Or.Functs:
				a.dataOut <- data1 ^ data2

			// srl - shift right logical
			case isa.Srl.Functs:
				a.dataOut <- data1 >> data2

			// sra - shift right arithmetic
			case isa.Sra.Functs:
				a.dataOut <- uint32(int32(data1) >> data2)

			// or
			case isa.Or.Functs:
				a.dataOut <- data1 | data2

			// and
			case isa.And.Functs:
				a.dataOut <- data1 & data2

			// mul
			case isa.Mul.Functs:
				a.dataOut <- data1 * data2
			case isa.Mulh.Functs:
				a.dataOut <- mulh(data1, data2)
			case isa.Mulhsu.Functs:
				a.dataOut <- mulhsu(data1, data2)
			case isa.Mulhu.Functs:
				a.dataOut <- mulhu(data1, data2)

			case isa.Div.Functs:
			case isa.Divu.Functs:
			case isa.Rem.Functs:
			case isa.Remu.Functs:
			}

		}
	}()

	return nil
}

func (a *ALU) Port() device.Port {
	return device.Port{
		Wires.Data1In: a.data1In,
		Wires.Data2In: a.data2In,
		Wires.DataOut: a.dataOut,
	}
}

func (a *ALU) Data1In(data device.WiresIn) {
	a.data1In = data
}

func (a *ALU) Data2In(data device.WiresIn) {
	a.data2In = data
}

func (a *ALU) DataOut() device.WiresOut {
	return a.dataOut
}

// FunctsIn this wire receives a 32-bit value in which
// the lower 10 bits encondes ISA fields funct7 and funct3:
//
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
func (a *ALU) FunctsIn(functs device.WiresIn) {
	a.functsIn = functs
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
