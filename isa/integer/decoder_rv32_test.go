//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package integer

import (
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

// TestDecoder tests R-format integer instructions
//
// fn7     RD2   RD1   fn3 RS    OPCODE
// 0000000_00010_00001_000_00101_0110011
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
			fields: datapath.OpFields{Shift: 0b00010,  Rs1: 0b00001, Funct3: 0b001, Rd: 0b00101, Opcode: 0b0010011},
		},
		{
			name:   Slti.Name,
			inst:   0b000000000010_00001_010_10101_0010011,
			fields: datapath.OpFields{Imm: 0b000000000010,  Rs1: 0b00001, Funct3: 0b010, Rd: 0b10101, Opcode: 0b0010011},
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
				t.Errorf("unexpected Op field value %08b", fields.Funct3)
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
				switch fields.Funct3{
				case Slli.F3, Srli.F3, Srai.F3:
					if fields.Shift != test.fields.Shift {
						t.Errorf("unexpected Shift field value %08b", fields.Shift)
					}
					if fields.Funct7 != test.fields.Funct7 {
						t.Errorf("unexpected Funct7 field value %08b", fields.Funct7)
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
