//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package store

import (
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name   string
		inst   datapath.XWord
		fields datapath.OpFields
	}{
		{
			name:   Sb.Name,
			inst:   0b0000100_10011_00001_000_00001_0100011,
			fields: datapath.OpFields{Rs2: 0b10011, Rs1: 0b00001, Funct3: 0b000, Imm: 0b000010000001},
		},
		{
			name:   Sh.Name,
			inst:   0b0010100_11011_01001_001_00101_0100011,
			fields: datapath.OpFields{Rs2: 0b11011, Rs1: 0b01001, Funct3: 0b001, Imm: 0b001010000101},
		},
		{
			name:   Sw.Name,
			inst:   0b0010100_11011_01001_010_00101_0100011,
			fields: datapath.OpFields{Rs2: 0b11011, Rs1: 0b01001, Funct3: 0b010, Imm: 0b001010000101},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.inst)
			if fields.Opcode != isa.Opcodes.S {
				t.Errorf("Unexpected Opcode %05b for op %s: %#v", fields.Opcode, test.name, fields)
			}
			if fields.Funct3 != test.fields.Funct3 {
				t.Errorf("Operation %s has unexpected Op %b", test.name, fields.Funct3)
			}
			if fields.Rs1 != test.fields.Rs1 {
				t.Errorf("Operation %s has unexpected RS1 %b", test.name, fields.Rs1)
			}
			if fields.Rs2 != test.fields.Rs2 {
				t.Errorf("Operation %s has unexpected RS2 %b", test.name, fields.Rs2)
			}
			if fields.Imm != test.fields.Imm {
				t.Errorf("Operation %s has unexpected Imm %b", test.name, fields.Imm)
			}
		})
	}
}
