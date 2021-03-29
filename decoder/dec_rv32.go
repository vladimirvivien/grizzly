package decoder

import (
	"encoding/binary"
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
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
	out   chan []byte
}

func New() *Decoder {
	dec := &Decoder{
		BaseComponent: datapath.NewBase(),
		out: make(chan []byte),
	}
	dec.Connect(Labels.OutFields, dec.out)
	return dec
}

// Run starts the decoder
func (d *Decoder) Run() error {
	instructions := d.GetPin(Labels.Instruction)
	if instructions == nil {
		return fmt.Errorf("decoder: input not set")
	}

	// launch main loop
	go func() {
		defer close(d.out)

		for {
			bits, opened := <-instructions
			if !opened {
				return
			}

			inst := instFromStream(bits)
			opcode := isa.GetOpcode(inst)

			var fields datapath.OpFields
			switch opcode {
			case isa.Opcodes.R, isa.Opcodes.RI:
				fields = integer.Decode(inst)
			case isa.Opcodes.L:
				fields = load.Decode(inst)
			case isa.Opcodes.S:
				fields = store.Decode(inst)
			}

			d.out <- datapath.EncodeOpFields(fields)
		}

	}()

	return nil
}

// decodeFromStream decodes input from stream to
// instruction Word
func instFromStream(bits []byte) datapath.XWord {
	return binary.LittleEndian.Uint32(bits)
}