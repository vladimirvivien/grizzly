package ctrlunit

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/reg"
)

// TestCtrl_Regfile tests controller/register file interaction
func TestCtrl_Regfile(t *testing.T) {
	ctrl := newCtrl()

	insts := datapath.MakeWires()
	ctrl.SetPin(In.Insts, insts)
	ctrl.SetClock(clock.New(2 * time.Millisecond))

	// connect controller and register
	regfile := reg.New().(*reg.RegisterFile)
	regData := datapath.MakeWires()
	regfile.SetPin(reg.In.Werf,ctrl.GetPin(Out.Werf))
	regfile.SetPin(reg.In.RS1Addr, ctrl.GetPin(Out.RS1))
	regfile.SetPin(reg.In.RS2Addr, ctrl.GetPin(Out.RS2))
	regfile.SetPin(reg.In.RDAddr, ctrl.GetPin(Out.RD))
	regfile.SetPin(reg.In.Data, regData)
	mux := device.Mux("alu-op", ctrl.GetPin(Out.ALUSrc), regfile.GetPin(reg.Out.RS2Data), ctrl.GetPin(Out.Imm))

	// prepare register
	regfile.SideLoad(2, 4)
	regfile.SideLoad(6, 12)
	regfile.SideLoad(8, 16)

	// R/I instructions
	go func() {
		insts <- 0b0000000_00110_00010_000_00101_0110011 // add  reg[5]  = reg[2]=4, reg[6]=12
		regData <- 4 + 12 // alu result
		insts <- 0b000000000010_00101_000_00100_0010011  // addi reg[4]  = reg[5]=16, 2
		regData <- 16+2 // alu result
		insts <- 0b0000000_01000_00110_000_00111_0110011 // add  reg[7]  = reg[6]=12, reg[8]=16
		regData <- 12+16 // alu result
		insts <- 0b0000000_00010_01000_000_01010_0110011 // add  reg[10] = reg[8]=16, reg[2]=4
		regData <- 16+4 // alu result
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
	memOp := ctrl.GetPin(Out.MemOp)
	memRen := ctrl.GetPin(Out.MemRen)
	wbSel := ctrl.GetPin(Out.WBSel)

	rcvr := datapath.NewReceiver("test:ctrl:reg")

	// add  reg[5]  = reg[2]=4, reg[6]=12
	data := rcvr.R(aluOp, rd1, mux, memOp, memRen, wbSel)
	aluInput1 := data[1]
	aluInput2 := data[2] // reg-alu mux output
	if aluInput1 != 4 {
		t.Fatalf("unexpected reg data1 %d", aluInput1)
	}
	if aluInput2 != 12 {
		t.Fatalf("unexpected mux data %d", aluInput2)
	}

	// addi reg[4]  = reg[5]=16, 2
	data = rcvr.R(aluOp, rd1, mux, memOp, memRen, wbSel)
	aluInput1 = data[1]
	aluInput2 = data[2]
	if aluInput1 != 16 {
		t.Fatalf("unexpected reg data1 %d", aluInput1)
	}
	if aluInput2 != 2 {
		t.Fatalf("unexpected mux data %d", aluInput2)
	}

	// add  reg[7]  = reg[6]=12, reg[8]=16
	data = rcvr.R(aluOp, rd1, mux, memOp, memRen, wbSel)
	aluInput1 = data[1]
	aluInput2 = data[2]
	if aluInput1 != 12 {
		t.Fatalf("unexpected reg data1 %d", aluInput1)
	}
	if aluInput2 != 16 {
		t.Fatalf("unexpected mux data %d", aluInput2)
	}

	// add  reg[10] = reg[8]=16, reg[2]=4
	data = rcvr.R(aluOp, rd1, mux, memOp, memRen, wbSel)
	aluInput1 = data[1]
	aluInput2 = data[2]
	if aluInput1 != 16 {
		t.Fatalf("unexpected reg data1 %d", aluInput1)
	}
	if aluInput2 != 4 {
		t.Fatalf("unexpected mux data %d", aluInput2)
	}
}
