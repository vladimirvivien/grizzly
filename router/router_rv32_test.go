package router

import (
	"math"
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

func TestRouter_Run(t *testing.T) {
	resultCh := make(chan []byte)
	go func (){
		// R-Opcode
		resultCh <- datapath.EncodeAluResult(datapath.AluResult{Opcode: isa.Opcodes.R, Funct3: integer.Add.F3, Funct7: integer.Add.F7, Rd: 0b00101, Value: math.MaxInt16})
		// RI-Opcode
		resultCh <- datapath.EncodeAluResult(datapath.AluResult{Opcode: isa.Opcodes.RI, Funct3: integer.Addi.F3, Funct7: integer.Addi.F7, Rd: 0b00101, Value: math.MaxInt8})
		close(resultCh)
	}()

	router := New()
	router.Connect(Labels.InAluResult,resultCh)
	if err := router.Run(); err != nil {
		t.Fatal(err)
	}

	// R-code
	output := <- router.GetPin(Labels.OutRegisterData)
	aluData := datapath.DecodeRegisterData(output)
	if aluData.Rd != 5{
		t.Errorf("unexpected reg addr: %d", aluData.Rd)
	}
	if aluData.Value != math.MaxInt16 {
		t.Errorf("unexpected reg data value %d", aluData.Value)
	}

	// RI-code
	output = <- router.GetPin(Labels.OutRegisterData)
	aluData = datapath.DecodeRegisterData(output)
	if aluData.Rd != 5{
		t.Errorf("unexpected reg addr: %d", aluData.Rd)
	}
	if aluData.Value != math.MaxInt8 {
		t.Errorf("unexpected reg data value %d", aluData.Value)
	}
}
