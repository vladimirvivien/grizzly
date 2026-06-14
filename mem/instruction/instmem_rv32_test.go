//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package instruction

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestInstructionMemory_New(t *testing.T) {
	mem := New(1024)
	if mem.outInst == nil {
		t.Errorf("missing output channel")
	}
	if mem.GetPin(Labels.OutInstruction) == nil {
		t.Errorf("pin not set %s", Labels.OutInstruction)
	}
}


// File
//main:
//addi x20, x0, 2
//addi x21, x0, 4
//addi x22, x0, 12
//addi x24, x20, 4
//add  x25, x21, x22
//slli x26, x21, 1
func TestInstructionMemory_LoadFile(t *testing.T) {
	tests := []struct {
		name string
		file string
		length int
		insts []datapath.XWord
	}{
		{
			name: "add.bin",
			file: "../../testing/programs/rtypes_rv32/add.bin",
			length:24, // total bytes
			insts: []datapath.XWord{
				0x00200a13,
				0x00400a93,
				0x00c00b13,
				0x004a0c13,
				0x016a8cb3,
				0x001a9d13,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			pcCh := make(chan []byte)
			go func(){
				for x := 0; x < test.length; x += 4 {
					pcCh <- datapath.EncodePC(datapath.ProgramCounter(x))
				}
				close(pcCh)
			}()

			mem, err := NewFromFile(test.file)
			if err != nil {
				t.Error(err)
			}
			mem.Connect(Labels.InPC, pcCh)

			if err := mem.Run(); err != nil {
				t.Error(err)
			}

			waiter := make(chan struct{})
			go func() {
				instCh := mem.GetPin(Labels.OutInstruction)
				var i int
				for stream := range instCh {
					inst := datapath.DecodeInstruction(stream)
					if test.insts[i] != inst.Inst {
						t.Errorf("unexpected instruction %d: %d", i, inst.Inst)
					}
					i++
				}
				close(waiter)
			}()

			select {
			case <-waiter:
			case <-time.After(50 * time.Millisecond):
				t.Fatal("DataMemory operations took too long to complete")
			}

		})
	}
}