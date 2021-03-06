package decoder

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

type Bytestream = <-chan []byte
type Decoder struct {
	bits  Bytestream
	clock clock.Clock
	width int
	out   chan datapath.OpFields
}

func New() *Decoder {
	return &Decoder{
		out:   make(chan datapath.OpFields, 0),
	}
}

func (d *Decoder) SetClock(c clock.Clock) {
	d.clock = c
}
func (d *Decoder) Input(in Bytestream) {
	d.bits = in
}
func (d *Decoder) SetInstructionWidth(w int) {
	d.width = w
}

func (d *Decoder) Output() <-chan datapath.OpFields {
	return d.out
}

func (d *Decoder) Run() error {
	if d.clock == nil {
		return fmt.Errorf("clock not set")
	}

	switch d.width {
	case datapath.Width32:
	default:
		return fmt.Errorf("unsupported insruction size: %d", d.width)
	}

	// launch main loop
	go func() {
		defer close(d.out)

		for {
			select {
			case _, opened := <-d.clock.Ticks():
				if !opened {
					return
				}

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
		}
	}()

	return nil
}
