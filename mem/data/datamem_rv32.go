package mem

import (
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

var(
	Labels = struct{
		InOperation datapath.Pin
		OutRegData  datapath.Pin
	}{
		InOperation: datapath.Pin("mem.in.operation"),
		OutRegData:  datapath.Pin("mem.out.reg_data"),
	}
)
type DataMemory struct {
	*datapath.BaseComponent
	*BaseMemory
	sync.RWMutex
	outReg chan []byte
}

func New(size uint64) *DataMemory {
	mem := &DataMemory{
		BaseComponent: datapath.NewBase(),
		BaseMemory: NewBase(size),
		outReg: make(chan []byte),
	}
	mem.Connect(Labels.OutRegData, mem.outReg)
	return  mem
}

func (m *DataMemory) Run() error {
	ops := m.GetPin(Labels.InOperation)
	if ops == nil {
		return fmt.Errorf("memory: missing input: %s", Labels.InOperation)
	}

	go func() {
		defer close(m.outReg)
		for {
			stream, opened := <-ops
			if !opened {
				return
			}

			op := datapath.DecodeMemOp(stream)
			switch op.Opcode {
			case isa.Opcodes.L:
				data := m.Read(op.Addr, op.Op)
				m.outReg <- datapath.EncodeRegData(datapath.RegisterData{
					Rd:    op.Rd,
					Value: data,
				})
			case isa.Opcodes.S:
				m.Write(op.Addr,op.Data, op.Op)
			}
		}
	}()

	return nil
}

func (m *DataMemory) TestSideLoad(addr datapath.XWord, val datapath.XWord) {
	m.Write(addr, val, store.Sw.F3)
}

func (m *DataMemory) TestProbe(addr datapath.XWord) datapath.XWord {
	return m.Read(addr, load.Lw.F3)
}