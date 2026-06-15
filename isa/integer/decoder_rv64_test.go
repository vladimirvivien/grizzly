//go:build rv64 || rv64i

package integer

import (
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

func TestDecoder(t *testing.T) {
	tests := []struct {
		name   string
		inst   datapath.XWord
		fields datapath.OpFields
	}{
		{
			name:   Add.Name,
			inst:   0b0000000_00010_00001_000_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0000000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b000, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   Sub.Name,
			inst:   0b0100000_00010_00001_000_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0100000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b000, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   Sll.Name,
			inst:   0b0000000_00010_00001_001_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0000000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b001, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   Slt.Name,
			inst:   0b0000000_00010_00001_010_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0000000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b010, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   Srl.Name,
			inst:   0b0000000_00010_00001_101_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0000000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b101, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   Sra.Name,
			inst:   0b0100000_00010_00001_101_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0100000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b101, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   Or.Name,
			inst:   0b0000000_00010_00001_110_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0000000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b110, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   And.Name,
			inst:   0b0000000_00010_00001_111_00101_0110011,
			fields: datapath.OpFields{Funct7: 0b0000000, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b111, Rd: 0b00101, Opcode: 0b0110011},
		},
		{
			name:   Addi.Name,
			inst:   0b000000000010_00001_000_00101_0010011,
			fields: datapath.OpFields{Imm: 0b000000000010, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b000, Rd: 0b00101, Opcode: 0b0010011},
		},
		{
			name:   Slli.Name,
			inst:   0b0000000_00010_00001_001_00101_0010011,
			fields: datapath.OpFields{Shift: 0b00010, Rs1: 0b00001, Funct3: 0b001, Rd: 0b00101, Opcode: 0b0010011},
		},
		{
			name:   Slti.Name,
			inst:   0b000000000010_00001_010_10101_0010011,
			fields: datapath.OpFields{Imm: 0b000000000010, Rs1: 0b00001, Funct3: 0b010, Rd: 0b10101, Opcode: 0b0010011},
		},
		{
			name:   Srli.Name,
			inst:   0b0000000_01010_00001_101_01101_0010011,
			fields: datapath.OpFields{Shift: 0b01010, Rs1: 0b00001, Funct3: 0b101, Rd: 0b01101, Opcode: 0b0010011},
		},
		{
			name:   Srai.Name,
			inst:   0b0100000_11011_01101_101_00101_0010011,
			fields: datapath.OpFields{Shift: 0b11011, Funct7: 0b0100000, Rs2: 0b00010, Rs1: 0b01101, Funct3: 0b101, Rd: 0b00101, Opcode: 0b0010011},
		},
		{
			name:   Ori.Name,
			inst:   0b010001000010_00001_110_00101_0010011,
			fields: datapath.OpFields{Imm: 0b010001000010, Rs1: 0b00001, Funct3: 0b110, Rd: 0b00101, Opcode: 0b0010011},
		},
		{
			name:   Andi.Name,
			inst:   0b010000000010_10101_111_10101_0010011,
			fields: datapath.OpFields{Imm: 0b010000000010, Rs1: 0b10101, Funct3: 0b111, Rd: 0b10101, Opcode: 0b0010011},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.inst)
			if fields.Opcode != test.fields.Opcode {
				t.Errorf("unexpected opcode %08b", fields.Opcode)
			}
			if fields.Rd != test.fields.Rd {
				t.Errorf("unexpected RD field value %08b", fields.Rd)
			}
			if fields.Funct3 != test.fields.Funct3 {
				t.Errorf("unexpected Funct3 field value %08b", fields.Funct3)
			}
			if fields.Rs1 != test.fields.Rs1 {
				t.Errorf("unexpected RS1 field value %08b", fields.Rs1)
			}
			switch fields.Opcode {
			case isa.Opcodes.R:
				if fields.Rs2 != test.fields.Rs2 {
					t.Errorf("unexpected RS2 field value %08b", fields.Rs2)
				}
				if fields.Funct7 != test.fields.Funct7 {
					t.Errorf("unexpected Funct7 field value %08b", fields.Funct7)
				}
				if fields.Imm != 0 {
					t.Errorf("unexpected Imm field value %08b", fields.Imm)
				}
				if fields.Shift != 0 {
					t.Errorf("unexpected Shift field value %08b", fields.Shift)
				}
			case isa.Opcodes.RI:
				switch fields.Funct3 {
				case Slli.F3, Srli.F3, Srai.F3:
					// In RV64, shift is 6 bits, Funct7 is 6 bits (shifted by 26)
					// Verify that shift is decoded correctly
					if fields.Shift != test.fields.Shift {
						t.Errorf("unexpected Shift field value %d, expected %d", fields.Shift, test.fields.Shift)
					}
				default:
					if fields.Imm != test.fields.Imm {
						t.Errorf("unexpected IMM field value %032b", fields.Imm)
					}
					if fields.Rs2 != 0 {
						t.Errorf("unexpected RS2 field value %08b", fields.Rs2)
					}
				}
			}
		})
	}
}

func FuzzDecodeInteger(f *testing.F) {
	f.Add(uint32(0b0000000_00010_00001_000_00101_0110011))
	f.Add(uint32(0b000000000010_00001_000_00101_0010011))
	f.Fuzz(func(t *testing.T, instVal uint32) {
		inst := datapath.XWord(instVal)
		opcode := inst & 0x7F
		if opcode != 0x33 && opcode != 0x13 {
			return
		}

		fields := Decode(inst)

		if fields.Opcode != uint8(opcode) {
			t.Errorf("opcode mismatch: got %x, expected %x", fields.Opcode, opcode)
		}
		if fields.Rd > 31 || fields.Rs1 > 31 {
			t.Errorf("registers must be <= 31: rd=%d, rs1=%d", fields.Rd, fields.Rs1)
		}
		if fields.Funct3 > 7 {
			t.Errorf("funct3 must be <= 7: %d", fields.Funct3)
		}

		if opcode == 0x33 {
			rs2 := uint8((inst >> 20) & 0x1F)
			f7 := uint8((inst >> 25) & 0x7F)
			if fields.Rs2 != rs2 {
				t.Errorf("rs2 mismatch: got %d, expected %d", fields.Rs2, rs2)
			}
			if fields.Funct7 != f7 {
				t.Errorf("funct7 mismatch: got %d, expected %d", fields.Funct7, f7)
			}
		} else if opcode == 0x13 {
			f3 := uint8((inst >> 12) & 0x7)
			switch f3 {
			case Slli.F3, Srli.F3, Srai.F3:
				shamt := uint8((inst >> 20) & 0x3F) // 6 bits for RV64
				f7 := uint8((inst >> 26) & 0x3F)
				if fields.Shift != shamt {
					t.Errorf("shift mismatch: got %d, expected %d", fields.Shift, shamt)
				}
				if fields.Funct7 != f7 {
					t.Errorf("funct7 mismatch: got %d, expected %d", fields.Funct7, f7)
				}
			default:
				imm := (inst >> 20) & 0xFFF
				if (imm & 0x800) != 0 {
					imm |= 0xfffff000
				}
				if fields.Imm != uint32(imm) {
					t.Errorf("imm mismatch: got %d, expected %d", fields.Imm, imm)
				}
			}
		}
	})
}
