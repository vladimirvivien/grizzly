package decoder

import (
	"encoding/binary"
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

func TestDecoder_Run(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*testing.T) datapath.Bytestream
		assess func(*testing.T, *Decoder)
	}{
		{
			name: "multi insts",
			setup: func(t *testing.T) datapath.Bytestream {
				stream := make(chan []byte)
				go func() {
					// R
					inst := make([]byte, 4)
					binary.LittleEndian.PutUint32(inst, 0b0000000_00010_00001_000_00101_0110011)
					stream <- inst

					// RI
					inst = make([]byte, 4)
					binary.LittleEndian.PutUint32(inst, 0b010001000010_00001_110_00101_0010011)
					stream <- inst

					// Load
					inst = make([]byte, 4)
					binary.LittleEndian.PutUint32(inst, 0b100010001011_01011_101_00101_0000011)
					stream <- inst

					// Store
					inst = make([]byte, 4)
					binary.LittleEndian.PutUint32(inst, 0b0010100_11011_01001_010_00101_0100011)
					stream <- inst

					close(stream)
				}()
				return stream
			},
			assess: func(t *testing.T, dec *Decoder) {
				// R
				fields := <-dec.Output()
				if fields.Opcode != isa.Opcodes.R {
					t.Errorf("unexpected field value %v", fields.Opcode)
				}
				if fields.Funct3 != 0 && fields.Funct7 != 0{
					t.Errorf("unexpected functs value %d, %d", fields.Funct7, fields.Funct3)
				}

				// RI
				fields = <-dec.Output()
				if fields.Opcode != isa.Opcodes.RI {
					t.Errorf("unexpected field value %v", fields.Opcode)
				}
				if fields.Imm != 0b010001000010{
					t.Errorf("unexpected imm value %d", fields.Imm)
				}

				// Load
				fields = <-dec.Output()
				if fields.Opcode != isa.Opcodes.L {
					t.Errorf("unexpected field value %v", fields.Opcode)
				}
				if fields.Imm != 0b100010001011{
					t.Errorf("unexpected imm value %d", fields.Imm)
				}
				if fields.Funct3 != 0b101{
					t.Errorf("unexpected Funct3 value %d", fields.Funct3)
				}

				// Store
				fields = <-dec.Output()
				if fields.Opcode != isa.Opcodes.S {
					t.Errorf("unexpected field value %v", fields.Opcode)
				}
				if fields.Imm != 0b001010000101{
					t.Errorf("unexpected imm value %d", fields.Imm)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			dec := New()
			dec.Input(test.setup(t))

			if err := dec.Run(); err != nil {
				t.Fatal(err)
			}

			test.assess(t, dec)
		})
	}
}
