package ctrlunit

import (
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

// TestCtrl_Regfile tests controller/register file interaction
func TestCtrl_Regfile(t *testing.T) {
	ctrl := newCtrl()

	insts := datapath.MakeWires()
	ctrl.SetPin(In.Insts, insts)

	// connect controller and register
	regfile := reg.New().(*reg.RegisterFile)
	regData := datapath.MakeWires()
	regfile.SetPin(reg.In.Werf,ctrl.GetPin(Out.Werf))
	regfile.SetPin(reg.In.RS1Addr, ctrl.GetPin(Out.RS1))
	regfile.SetPin(reg.In.RS2Addr, ctrl.GetPin(Out.RS2))
	regfile.SetPin(reg.In.RDAddr, ctrl.GetPin(Out.RD))
	regfile.SetPin(reg.In.Data, regData)
	mux := device.Mux(ctrl.GetPin(Out.ALUSrc), regfile.GetPin(reg.Out.RS2Data), ctrl.GetPin(Out.Imm))

	// prepare register
	regfile.SideLoad(2, 4)
	regfile.SideLoad(6, 12)
	regfile.SideLoad(8, 16)

	// R/I instructions
	go func() {
		insts <- 0b0000000_00110_00010_000_00101_0110011 // add  reg[5]  = reg[2]=4, reg[6]=12
		regData <- 4 + 12
		insts <- 0b000000000010_00101_000_00101_0010011  // addi reg[5]  = reg[5]=16, 2
		regData <- 5 + 18
		insts <- 0b0000000_01000_00110_000_00111_0110011 // add  reg[7]  = reg[6]=12, reg[8]=16
		regData <- 12+16
		insts <- 0b0000000_00010_01000_000_01010_0110011 // add  reg[10] = reg[8]=16, reg[2]=4
		regData <- 16+4
	}()

	// start components
	if err := regfile.Run(); err != nil {
		t.Fatal(err)
	}
	if err := ctrl.Run(); err != nil {
		t.Fatal(err)
	}

	// connect to output pins and test values
	aluOp := ctrl.GetPin(Out.ALUOp)
	rd1 := regfile.GetPin(reg.Out.RS1Data)

	// add  reg[5]  = reg[2]=4, reg[6]=12
	data := datapath.Collect(aluOp, rd1, mux)
	regData1 := data[1]
	muxData := data[2]
	if regData1 != 4 {
		t.Fatalf("unexpected reg data1 %d", regData1)
	}
	if muxData != 12 {
		t.Fatalf("unexpected mux data %d", muxData)
	}

	// addi reg[5]  = reg[5]=16, 2
	data = datapath.Collect(aluOp, rd1, mux)
	regData1 = data[1]
	muxData = data[2]
	if regData1 != 16 {
		t.Fatalf("unexpected reg data1 %d", regData1)
	}
	if muxData != 2 {
		t.Fatalf("unexpected mux data %d", muxData)
	}

	// add  reg[7]  = reg[6]=12, reg[8]=16
	data = datapath.Collect(aluOp, rd1, mux)
	regData1 = data[1]
	muxData = data[2]
	if regData1 != 12 {
		t.Fatalf("unexpected reg data1 %d", regData1)
	}
	if muxData != 16 {
		t.Fatalf("unexpected mux data %d", muxData)
	}

	// add  reg[10] = reg[8]=16, reg[2]=4
	data = datapath.Collect(aluOp, rd1, mux)
	regData1 = data[1]
	muxData = data[2]
	if regData1 != 16 {
		t.Fatalf("unexpected reg data1 %d", regData1)
	}
	if muxData != 4 {
		t.Fatalf("unexpected mux data %d", muxData)
	}
}
