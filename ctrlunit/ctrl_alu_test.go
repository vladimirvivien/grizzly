package ctrlunit

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
)

func TestCtrl_ALU(t *testing.T) {
	ctrl := newCtrl()
	insts := datapath.MakeWires()
	ctrl.SetPin(In.Insts, insts)
	ctrl.SetClock(clock.New(2 * time.Millisecond))

	op1 := datapath.MakeWires()
	rd2 := datapath.MakeWires() // fake reg data2 for testing
	arlou := alu.New().(*alu.ALU)
	arlou.SetPin(alu.In.Operation, ctrl.GetPin(Out.ALUOp))
	arlou.SetPin(alu.In.Operand1, op1)
	arlou.SetPin(alu.In.Operand2, device.Mux(ctrl.GetPin(Out.ALUSrc), rd2, ctrl.GetPin(Out.Imm)))
	// write back mux use WBSel line to route from either alu or mem data (simulated)
	memData := datapath.MakeWires()
	wbMux := device.Mux(ctrl.GetPin(Out.WBSel), arlou.GetPin(alu.Out.Result), memData)



	go func() {
		insts <- 0b0000000_00110_00010_000_00101_0110011 // add  reg[5]  = reg[2]=4, reg[6]=12
		op1 <- 4; rd2 <- 12; memData <- 0 // simulate alu-operand-1; reg data out 2; mem-addr out
		insts <- 0b000000000010_00101_000_00101_0010011  // addi reg[5]  = reg[5]=16, 2
		op1 <- 16; rd2 <- 0; memData <- 0
		insts <- 0b0000000_01000_00110_000_00111_0110011 // add  reg[7]  = reg[6]=12, reg[8]=16
		op1 <- 12; rd2 <-16; memData <- 0
		insts <- 0b0000000_00010_01000_000_01010_0110011 // add  reg[10] = reg[8]=16, reg[2]=4
		op1 <- 16; rd2 <- 4; memData <- 0
	}()

	if err := arlou.Run(); err != nil {
		t.Fatal(err)
	}
	if err := ctrl.Run(); err != nil {
		t.Fatal(err)
	}

	rs1, rs2, memRead, werf, rd :=
		ctrl.GetPin(Out.RS1),
		ctrl.GetPin(Out.RS2),
		ctrl.GetPin(Out.MemRead),
		ctrl.GetPin(Out.Werf),
		ctrl.GetPin(Out.RD)

	// process inst 1, add
	datapath.Collect(rs1, rs2, memRead, werf, rd)
	if result := datapath.Collect(wbMux)[0]; result != 16 {
		t.Fatal("unexpected alu result:", result)
	}

	// process inst 2, addi
	datapath.Collect(rs1, rs2, memRead, werf, rd)
	if result := datapath.Collect(wbMux)[0]; result != 18 {
		t.Fatal("unexpected alu result:", result)
	}

	// process inst 3, add
	datapath.Collect(rs1, rs2, memRead, werf, rd)
	if result := datapath.Collect(wbMux)[0]; result != 28 {
		t.Fatal("unexpected alu result:", result)
	}

	// process inst 4, add
	datapath.Collect(rs1, rs2, memRead, werf, rd)
	if result := datapath.Collect(wbMux)[0]; result != 20 {
		t.Fatal("unexpected alu result:", result)
	}
}
