//go:build rv64 || rv64i

package instruction

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/mem"
)

var (
	Labels = struct {
		InPC           datapath.Pin
		OutInstruction datapath.Pin
	}{
		InPC:           datapath.Pin("instmem.in.program_counter"),
		OutInstruction: datapath.Pin("mem.out.instruction"),
	}
)

type InstructionMemory struct {
	*datapath.BaseComponent
	*mem.BaseMemory
	sync.RWMutex
	outInst chan []byte
}

func New(size uint64) *InstructionMemory {
	im := &InstructionMemory{
		BaseComponent: datapath.NewBase(),
		BaseMemory:    mem.NewBase(size),
		outInst:       make(chan []byte),
	}
	im.Connect(Labels.OutInstruction, im.outInst)
	return im
}

func NewFromFile(path string) (*InstructionMemory, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}
	m := New(uint64(len(file)))
	m.SetStore(file)
	return m, nil
}

func (m *InstructionMemory) Run() error {
	pcCh := m.GetPin(Labels.InPC)
	if pcCh == nil {
		return fmt.Errorf("inst memory: missing input: %s", Labels.InPC)
	}

	go func() {
		defer close(m.outInst)
		for {
			stream, opened := <-pcCh
			if !opened {
				return
			}

			pc := datapath.DecodePC(stream)
			var inst uint32
			if uint64(pc) >= uint64(m.GetSize()) {
				inst = 0x00000013 // NOP (addi x0, x0, 0)
			} else {
				inst = uint32(m.Read(pc, load.Lw.F3))
			}
			m.outInst <- datapath.EncodeInstruction(datapath.Instruction{
				PC:   pc,
				Inst: inst,
			})
		}
	}()

	return nil
}
