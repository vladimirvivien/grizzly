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

type Decoder struct {
	bits  datapath.Bytestream
	out   chan []byte
}

func New() *Decoder {
	return &Decoder{
		out: make(chan []byte, 0),
	}
}

// Input is connected to bytestream from instruction memory
func (d *Decoder) Input(in datapath.Bytestream) {
	d.bits = in
}

// Output returns a bytestream containing the
// decoded instruction fields laid out as:
//
// 0       1       2       3       4       5       6
// 01234567012345670123456701234567012345670123456701234567
// +-------+-------+-------+-------+-------+-------+------+
// |OpCode |   Rd  |Funct3 |   Rs1 |  Rs2  |Funct7 |Shift |
// +-------+-------+-------+-------+-------+-------+------+
// |               Imm             |
// +-------------------------------+
//
func (d *Decoder) Output() datapath.Bytestream {
	return d.out
}

// Run starts the decoder
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

			d.out <- encode(fields)
		}

	}()

	return nil
}

// encode encodes the decoded instruction fields into stream:
//
// 0       1       2       3       4       5       6
// 01234567012345670123456701234567012345670123456701234567
// +-------+-------+-------+-------+-------+-------+------+
// |OpCode |   Rd  |Funct3 |   Rs1 |  Rs2  |Funct7 |Shift |
// +-------+-------+-------+-------+-------+-------+------+
// |               Imm             |
// +-------------------------------+
//
func encode(f datapath.OpFields) []byte {
	buf := make([]byte,11,11)
	buf[0] = f.Opcode
	buf[1] = f.Rd
	buf[2] = f.Funct3
	buf[3] = f.Rs1
	buf[4] = f.Rs2
	buf[5] = f.Funct7
	buf[6] = f.Shift
	binary.LittleEndian.PutUint32(buf[7:], f.Imm)
	return buf
}