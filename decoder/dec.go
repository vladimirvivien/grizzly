package decoder

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

type Decoder struct {
	bits  datapath.Bytestream
	out   chan datapath.OpFields
}

func New() *Decoder {
	return &Decoder{
		out: make(chan datapath.OpFields, 0),
	}
}

func (d *Decoder) Input(in datapath.Bytestream) {
	d.bits = in
}

func (d *Decoder) Output() <-chan datapath.OpFields {
	return d.out
}

func (d *Decoder) Run() error {
	if d.bits == nil {
		return fmt.Errorf("decoder: input not set")
	}

	// launch main loop
	go func() {
		defer close(d.out)

		for {
			bits, opened := <-d.bits
			if !opened {
				return
			}

			inst := bytesToInst(bits)
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
			d.out <- fields
		}

	}()

	return nil
}
