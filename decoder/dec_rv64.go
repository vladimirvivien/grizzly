//go:build rv64 || rv64i

package decoder

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/branch"
	"github.com/vladimirvivien/grizzly/isa/integer"
	"github.com/vladimirvivien/grizzly/isa/jump"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

var (
	Labels = struct {
		Instruction datapath.Pin
		OutFields   datapath.Pin
	}{
		Instruction: datapath.Pin("decoder.in.instruction"),
		OutFields:   datapath.Pin("decoder.out.fields"),
	}
)

type Decoder struct {
	*datapath.BaseComponent
	out chan []byte
}

func New() *Decoder {
	dec := &Decoder{
		BaseComponent: datapath.NewBase(),
		out:           make(chan []byte),
	}
	dec.Connect(Labels.OutFields, dec.out)
	return dec
}

func (d *Decoder) Run() error {
	instructions := d.GetPin(Labels.Instruction)
	if instructions == nil {
		return fmt.Errorf("decoder: input not set")
	}

	go func() {
		defer close(d.out)

		for {
			bits, opened := <-instructions
			if !opened {
				return
			}

			inst := datapath.DecodeInstruction(bits)
			opcode := isa.GetOpcode(datapath.XWord(inst.Inst))

			var fields datapath.OpFields
			switch opcode {
			case isa.Opcodes.R, isa.Opcodes.RI:
				fields = integer.Decode(datapath.XWord(inst.Inst))
			case isa.Opcodes.L:
				fields = load.Decode(datapath.XWord(inst.Inst))
			case isa.Opcodes.S:
				fields = store.Decode(datapath.XWord(inst.Inst))
			case isa.Opcodes.J, isa.Opcodes.JI:
				fields = jump.Decode(datapath.XWord(inst.Inst))
			case isa.Opcodes.B:
				fields = branch.Decode(datapath.XWord(inst.Inst))
			}

			fields.PC = inst.PC
			d.out <- datapath.EncodeOpFields(fields)
		}
	}()

	return nil
}
