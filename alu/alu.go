package alu

import (
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/inst"
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
	data1In device.WiresIn
	data2In device.WiresIn
	dataOut device.Wires

	// functIn is a 32-bit value that
	// conctenates the bits from Funct7+Funct3
	// [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
	functIn device.WiresIn
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
			funct := <-a.functIn

			switch funct {
			case inst.Add.Funct:
				a.dataOut <- data1 + data2
			case inst.Sub.Funct:
				a.dataOut <- data1 - data2
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

// FunctIn receives a 32-bit value that is a
// concatenation of the bits from Funct7+Funct3
//
//    [XXXXXXXX XXXXXXXX XXXXXX77 77777333]
func (a *ALU) FunctIn(funct device.WiresIn) {
	a.functIn = funct
}
