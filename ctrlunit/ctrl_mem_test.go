package ctrlunit

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/mem"
)

func TestCtrl_Mem(t *testing.T){
	ctrl := newCtrl()
	insts := datapath.MakeWires()
	ctrl.SetPin(In.Insts, insts)
	ctrl.SetClock(clock.New(2 * time.Millisecond))

	aluOutput := datapath.MakeWires() // simulates ALU result
	// aluMemConnect splits ALU output into 1) a mem addr 2) or bypass to register
	aluMemConnect := device.Fanout(aluOutput, 2)

	memory := mem.New(1024).(*mem.Memory)
	memory.SetPin(mem.In.WriteEnable, ctrl.GetPin(Out.MemWrite))
	memory.SetPin(mem.In.ReadEnable, ctrl.GetPin(Out.MemRead))
	memory.SetPin(mem.In.Address, aluMemConnect[0])
	// wbMux: outputs the register write back value either from alu or from mem
	wbMux := device.Mux(ctrl.GetPin(Out.WBSel), aluMemConnect[1], memory.GetPin(mem.Out.DataRead))

	go func() {
		memory.TestSideLoad(1000, 24)
		insts <- 0b000010000000_00000_000_00101_0000011 // load
		aluOutput <- 1000 // effective addr

		memory.TestSideLoad(844, 2021)
		insts <- 0b000100000001_00000_000_00101_0000011 // load
		aluOutput <- 844 // effective addr

		insts <- 0b0000000_01000_00110_000_00111_0110011 // add
		aluOutput <- 12010 // alu output
	}()

	// start components
	if err := memory.Run(); err != nil {
		t.Fatal(err)
	}

	if err := ctrl.Run(); err != nil {
		t.Fatal(err)
	}

	// unconnected control wires
	aluOp, rs1, rs2, imm, aluSrc, werf, rd :=
		ctrl.GetPin(Out.ALUOp),
		ctrl.GetPin(Out.RS1),
		ctrl.GetPin(Out.RS2),
		ctrl.GetPin(Out.Imm),
   	    ctrl.GetPin(Out.ALUSrc),
		ctrl.GetPin(Out.Werf),
		ctrl.GetPin(Out.RD)

	// load mem[1000]
	data := datapath.Collect(aluOp, rs1, rs2, imm, aluSrc, werf, rd, wbMux)
	// check imm value 128
	if data[3] != 128 {
		t.Errorf("unexpected immediate: %012b", data[6])
	}
	// check mem data
	if data[7] != 24 {
		t.Errorf("unexpected memory data: %032b", data[7])
	}

	// load mem[844]
	data = datapath.Collect(aluOp, rs1, rs2, imm, aluSrc, werf, rd, wbMux)
	// check imm value 257
	if data[3] != 257 {
		t.Errorf("unexpected immediate: %012b", data[6])
	}
	// check mem data
	if data[7] != 2021 {
		t.Errorf("unexpected memory data: %032b", data[7])
	}

	//add
	data = datapath.Collect(aluOp, rs1, rs2, imm, aluSrc, werf, rd, wbMux)
	// check alu out == write back mux value
	if data[7] != 12010 {
		t.Errorf("unexpected write back value: %032b", data[7])
	}

}
