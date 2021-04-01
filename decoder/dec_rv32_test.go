package decoder

import (
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
					inst := datapath.EncodeInstruction(datapath.Instruction{PC: 0, Inst: 0b0000000_00010_00001_000_00101_0110011})
					stream <- inst

					// RI
					inst = datapath.EncodeInstruction(datapath.Instruction{PC: 4, Inst: 0b010001000010_00001_110_00101_0010011})
					stream <- inst

					// Load
					inst = datapath.EncodeInstruction(datapath.Instruction{PC: 8, Inst: 0b100010001011_01011_101_00101_0000011})
					stream <- inst

					// Store
					inst = datapath.EncodeInstruction(datapath.Instruction{PC: 12, Inst: 0b0010100_11011_01001_010_00101_0100011})
					stream <- inst

					close(stream)
				}()
				return stream
			},
			assess: func(t *testing.T, dec *Decoder) {
				// R
				fields := datapath.DecodeOpFields(<-dec.GetPin(Labels.OutFields))

				if fields.Opcode != isa.Opcodes.R {
					t.Errorf("unexpected opcode %v", fields.Opcode)
				}
				if fields.Funct3 != 0 && fields.Funct7 != 0 {
					t.Errorf("unexpected functs value %d, %d", fields.Funct3, fields.Funct7)
				}

				// RI
				fields = datapath.DecodeOpFields(<-dec.GetPin(Labels.OutFields))
				if fields.Opcode != isa.Opcodes.RI {
					t.Errorf("unexpected field value %v", fields.Opcode)
				}

				if fields.Imm != 0b010001000010 {
					t.Errorf("unexpected imm value %d", fields.Imm)
				}

				// Load
				fields = datapath.DecodeOpFields(<-dec.GetPin(Labels.OutFields))
				if fields.Opcode != isa.Opcodes.L {
					t.Errorf("unexpected field value %v", fields.Opcode)
				}
				if fields.Imm != 0b100010001011 {
					t.Errorf("unexpected imm value %d", fields.Imm)
				}
				if fields.Funct3 != 0b101 {
					t.Errorf("unexpected Op value %d", fields.Funct3)
				}

				// Store
				fields = datapath.DecodeOpFields(<-dec.GetPin(Labels.OutFields))
				if fields.Opcode != isa.Opcodes.S {
					t.Errorf("unexpected field value %v", fields.Opcode)
				}
				if fields.Imm != 0b001010000101 {
					t.Errorf("unexpected imm value %d", fields.Imm)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			dec := New()
			dec.Connect(Labels.Instruction, test.setup(t))

			if err := dec.Run(); err != nil {
				t.Fatal(err)
			}

			test.assess(t, dec)
		})
	}
}
