//go:build rv64 || rv64i

package core

import (
	"sync"
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/brancher"
	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/decoder"
	"github.com/vladimirvivien/grizzly/mem/data"
	"github.com/vladimirvivien/grizzly/mem/instruction"
	"github.com/vladimirvivien/grizzly/pc"
	"github.com/vladimirvivien/grizzly/reg"
)

type Core struct {
	clock    *clock.Clock
	pc       *pc.PC
	imem     *instruction.InstructionMemory
	dec      *decoder.Decoder
	reg      *reg.RegisterFile
	brancher *brancher.Brancher
	alu      *alu.ALU
	dmem     *data.DataMemory
}

func New() *Core {
	return &Core{
		clock:    clock.New(2 * time.Microsecond),
		pc:       pc.New(),
		imem:     instruction.New(1024),
		dec:      decoder.New(),
		reg:      reg.New(),
		brancher: brancher.New(),
		alu:      alu.New(),
		dmem:     data.New(1024 * 1000),
	}
}

func (c *Core) Run() error {
	c.wireAll()
	return c.startComponents()
}

func (c *Core) wireAll() {
	c.imem.Connect(instruction.Labels.InPC, c.pc.GetPin(pc.Labels.OutCounter))
	c.dec.Connect(decoder.Labels.Instruction, c.imem.GetPin(instruction.Labels.OutInstruction))
	c.reg.Connect(reg.Labels.InFields, c.dec.GetPin(decoder.Labels.OutFields))
	c.brancher.Connect(brancher.Labels.InBranchOp, c.reg.GetPin(reg.Labels.OutBranchOps))

	aluOps := merge(c.reg.GetPin(reg.Labels.OutAluOps), c.brancher.GetPin(brancher.Labels.OutOperation))
	c.alu.Connect(alu.Labels.InOperations, aluOps)

	c.reg.Connect(reg.Labels.InAluData, c.alu.GetPin(alu.Labels.OutRegData))
	c.dmem.Connect(data.Labels.InOperation, c.alu.GetPin(alu.Labels.OutMemOp))
	c.reg.Connect(reg.Labels.InMemData, c.dmem.GetPin(data.Labels.OutRegData))
	c.pc.Connect(pc.Labels.InPcOp, c.alu.GetPin(alu.Labels.OutPcOp))
}

func (c *Core) startComponents() error {
	components := []datapath.Component{c.pc, c.imem, c.dec, c.reg, c.brancher, c.alu, c.dmem}
	for _, comp := range components {
		if err := comp.Run(); err != nil {
			return err
		}
	}
	return nil
}

func merge(ch1, ch2 datapath.Bytestream) datapath.Bytestream {
	out := make(chan []byte)
	go func() {
		defer close(out)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for v := range ch1 {
				out <- v
			}
		}()
		go func() {
			defer wg.Done()
			for v := range ch2 {
				out <- v
			}
		}()
		wg.Wait()
	}()
	return out
}
