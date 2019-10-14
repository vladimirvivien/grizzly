package alu

import (
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/inst"
)

var (
	Wires = struct {
		Data1In,
		Data2In,
		FuncIn,
		DataOut string
	}{
		Data1In: "alu.data1.in",
		Data2In: "alu.data2.in",
		FuncIn:  "alu.func.in",
		DataOut: "alu.data.out",
	}
)

type ALUFunc struct {
	Funct3 uint32
	Funct7 uint32
}

type ALU struct {
	data1In device.WiresIn
	data2In device.WiresIn
	funcIn  <-chan ALUFunc
	dataOut device.Wires
}

func New() device.Type {
	return newAlu()
}

func newAlu() *ALU {
	return &ALU{
		dataOut: device.MakeWires(),
	}
}

func (a *ALU) Run() error {
	go func() {
		defer close(a.dataOut)
		for {
			op1 := <-a.data1In
			op2 := <-a.data2In
			op := <-a.funcIn

			switch op.Funct7 {
			case inst.Funct7Lo:
				switch op.Funct3 {
				case inst.Add.Funct3:
					a.dataOut <- op1 + op2
				}

			case inst.Funct7Hi:
				switch op.Funct3 {
				case inst.Sub.Funct3:
					a.dataOut <- op1 - op2
				}
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
