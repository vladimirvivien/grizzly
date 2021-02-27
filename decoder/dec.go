package decoder

import (
	"encoding/binary"
	"fmt"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

type Bitstream = <-chan []byte
type Decoder struct {
	bits Bitstream
	clock clock.Clock
	width int
}

func New() *Decoder {
	return &Decoder{width:datapath.XLEN}
}

func (d *Decoder) SetClock(c clock.Clock){
	d.clock = c
}
func (d *Decoder) SetBitstream(in Bitstream){
	d.bits = in
}
func (d *Decoder) SetInstructionWidth(w int){

}

func (d *Decoder) Run() error {
	if d.clock == nil {
		return fmt.Errorf("clock not set")
	}

	switch d.width{
	case datapath.Width32:
	default:
		return fmt.Errorf("unsupported insruction size: %d", d.width)
	}

	// launch main loop
	go func() {
		for range d.clock.Ticks() {
			bits := <- d.bits
			switch d.width{
			case datapath.Width32:
				inst := binary.LittleEndian.Uint32(bits)
				integer.Decode(inst)
			default:

			}
		}
	}()

	return nil
}