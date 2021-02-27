package alu

import (
	"log"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
)

var (
	Ops = struct {
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
		Operand1:  "alu.operand1.in",
		Operand2:  "alu.operand2.in",
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
	resultOut datapath.Wires // output
	zeroOut   datapath.Wires // zero line
	*device.Base
}

func New() device.Type {
	return newAlu()
}

func newAlu() *ALU {
	a := &ALU{
		resultOut: datapath.MakeWires(),
		Base:      device.NewBase(),
	}

	a.SetPin(Out.Result, a.resultOut)

	return a
}

// Run starts the ALU.
// Data1 and Data2 are read sequentially and must
// be available or risk blocking.
func (a *ALU) Run() error {
	log.Println("alu: starting...")
	go func() {
		defer close(a.resultOut)

		//rcvr := datapath.NewReceiver("alu")
		opWire := a.GetPin(In.Operation)
		operandWire1 := a.GetPin(In.Operand1)
		operandWire2 := a.GetPin(In.Operand2)

		for {
			// wait for opCtrl, receive input data
			op := <-opWire
			log.Printf("alu: aluOp %d", op)
			//data := rcvr.R(operandWire1, operandWire2)
			//data1, data2 := data[0], data[1]
			data1 := <-operandWire1
			log.Printf("alu: operand1 %d", data1)
			data2 := <-operandWire2
			log.Printf("alu: opreand2 %d", data2)

			log.Printf("alu: op=%05b, op1=%032b, op2=%032b", op, data1, data2)

			var opResult datapath.Word

			switch op {
			// addition: add, addi
			case Ops.Add:
				opResult = data1 + data2

			// sub
			case Ops.Sub:
				opResult = data1 - data2

			// shift logical left: sll, slli
			case Ops.Sll:
				opResult = data1 << data2

			// set if less then (signed): slt, slti
			case Ops.Slt:
				var result uint32
				if int32(data1) < int32(data2) {
					result = 1
				}
				opResult = result

			// set if less then (unsigned): sltu, sltiu
			case Ops.Sltu:
				var result uint32
				if data1 < data2 {
					result = 1
				}
				opResult = result

			// or, ori
			case Ops.Xor:
				opResult = data1 ^ data2

			// shift right logical: srl, srli
			case Ops.Srl:
				opResult = data1 >> data2

			// shift right arithmetic: sra, srai
			case Ops.Sra:
				opResult = uint32(int32(data1) >> data2)

			// or, ori
			case Ops.Or:
				opResult = data1 | data2

			// and, andi
			case Ops.And:
				opResult = data1 & data2

			// mul
			case Ops.Mul:
				opResult = data1 * data2
			case Ops.Mulh:
				opResult = mulh(data1, data2)
			case Ops.Mulhsu:
				opResult = mulhsu(data1, data2)
			case Ops.Mulhu:
				opResult = mulhu(data1, data2)

			case Ops.Div:
			case Ops.Divu:
			case Ops.Rem:
			case Ops.Remu:
			}

			a.resultOut <- opResult
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
