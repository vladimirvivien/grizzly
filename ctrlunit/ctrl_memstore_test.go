package ctrlunit

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
	"github.com/vladimirvivien/grizzly/mem"
)

func TestCtrl_MemStore(t *testing.T) {
	ctrl := newCtrl()
	insts := datapath.MakeWires()
	ctrl.SetPin(In.Insts, insts)
	ctrl.SetClock(clock.New(2 * time.Millisecond))

	// Rs2Data simulates reg[rs2] addr
	rs2Out := datapath.MakeWires()
	// regMemConnect splits reg[rs2] data 1) for alu (via mux) 2) to memory
	regMemConnect := device.Connector("reg-out-connector", rs2Out, 2)
	aluSrcMux := device.Mux("alu-op", ctrl.GetPin(Out.ALUSrc), regMemConnect[0], ctrl.GetPin(Out.Imm))

	// memAddr simulates generated addr from ALU
	memAddr := datapath.MakeWires()
	// aluMemConnect splits ALU output into 1) a mem addr 2) or bypass to register
	aluMemConnect := device.Connector("alu-out-connector",memAddr, 2)

	memory := mem.New(1024).(*mem.Memory)
	memory.SetPin(mem.In.WriteEnable, ctrl.GetPin(Out.MemWen))
	memory.SetPin(mem.In.Operation, ctrl.GetPin(Out.MemOp))
	memory.SetPin(mem.In.Address, aluMemConnect[0])
	memory.SetPin(mem.In.DataWrite, regMemConnect[1])

	// wbMux: outputs the register write back value either from alu or from mem
	wbMux := device.Mux("reg-wb", ctrl.GetPin(Out.WBSel), aluMemConnect[1], memory.GetPin(mem.Out.DataRead))

	go func() {
		// store byte
		insts <- 0b0000101_00100_00011_000_00001_0100011
		rs2Out <- 112 // mem data
		memAddr <- 24 // effective addr

		// store half word
		insts <- 0b0000111_10100_00111_001_10001_0100011
		rs2Out <- 23456 // mem data
		memAddr <- 1020 // effective addr

		// store half word
		insts <- 0b0000111_10100_00111_010_10001_0100011
		rs2Out <- 23456 // mem data
		memAddr <- 0000 // effective addr
	}()

	// start components
	if err := memory.Run(); err != nil {
		t.Fatal(err)
	}

	if err := ctrl.Run(); err != nil {
		t.Fatal(err)
	}

	// unconnected control wires
	aluOp, rs1, rs2, werf, rd :=
		ctrl.GetPin(Out.ALUOp),
		ctrl.GetPin(Out.RS1),
		ctrl.GetPin(Out.RS2),
		ctrl.GetPin(Out.Werf),
		ctrl.GetPin(Out.RD)

	tests := []struct {
		name string
		rs1,
		rs2,
		imm,
		writeBack,
		storedVal datapath.XWord
	}{
		{
			name:      "store byte (sb)",
			rs1:       0b00011,
			rs2:       0b00100,
			imm:       0b000010100001,
			writeBack: 24,
			storedVal: 112,
		},
		{
			name:      "store half word (sh)",
			rs1:       0b00111,
			rs2:       0b10100,
			imm:       0b000011110001,
			writeBack: 1020,
			storedVal: 23456,
		},
		{
			name:      "store word (sw)",
			rs1:       0b00111,
			rs2:       0b10100,
			imm:       0b000011110001,
			writeBack: 0,
			storedVal: 23456,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rcvr := datapath.NewReceiver("test:ctrl:memstore")
			// store value
			data := rcvr.R(aluOp, rs1, rs2, aluSrcMux, werf, rd, wbMux)
			// check rs1
			if data[1] != test.rs1 {
				t.Errorf("unexpected rs1: %012b", data[1])
			}
			// check rs2
			if data[2] != test.rs2 {
				t.Errorf("unexpected rs2: %012b", data[2])
			}
			// check imm value 128
			if data[3] != test.imm {
				t.Errorf("unexpected immediate: %012b", data[3])
			}
			// check wribeback mux: should output address
			if data[6] != test.writeBack {
				t.Errorf("unexpected memory data: %032b", data[6])
			}

			memVal := memory.TestProbe(test.writeBack)
			if memVal != test.storedVal {
				t.Errorf("unexpected memory value M[%032b]=%032b", 24, memVal)
			}
		})
	}
}
