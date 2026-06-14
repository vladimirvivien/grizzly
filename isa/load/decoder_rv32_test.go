//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package load

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
			name:   Lb.Name,
			inst:   0b000000010011_00001_000_00001_0000011,
			fields: datapath.OpFields{Imm: 0b000000010011, Rs1: 0b00001, Funct3: 0b000, Rd: 0b00001},
		},
		{
			name:   Lbu.Name,
			inst:   0b100000011011_01001_100_11001_0000011,
			fields: datapath.OpFields{Imm: 0xFFFFF81B, Rs1: 0b01001, Funct3: 0b100, Rd: 0b11001},
		},
		{
			name:   Lh.Name,
			inst:   0b000010001011_01011_001_11101_0000011,
			fields: datapath.OpFields{Imm: 0b000010001011, Rs1: 0b01011, Funct3: 0b001, Rd: 0b11101},
		},
		{
			name:   Lhu.Name,
			inst:   0b100010001011_01011_101_00101_0000011,
			fields: datapath.OpFields{Imm: 0xFFFFF88B, Rs1: 0b01011, Funct3: 0b101, Rd: 0b00101},
		},
		{
			name:   Lw.Name,
			inst:   0b111111001011_11011_110_00101_0000011,
			fields: datapath.OpFields{Imm: 0xFFFFFFCB, Rs1: 0b11011, Funct3: 0b110, Rd: 0b00101},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.inst)
			if fields.Opcode != isa.Opcodes.L {
				t.Errorf("Unexpected Opcode %05b for op %s: %#v", fields.Opcode, test.name, fields)
			}
			if fields.Funct3 != test.fields.Funct3 {
				t.Errorf("Operation %s has unexpected Op %b", test.name, fields.Funct3)
			}
			if fields.Rd != test.fields.Rd {
				t.Errorf("Operation %s has unexpected RD %b", test.name, fields.Rd)
			}
			if fields.Rs1 != test.fields.Rs1 {
				t.Errorf("Operation %s has unexpected RS1 %b", test.name, fields.Rs1)
			}
			if fields.Imm != test.fields.Imm {
				t.Errorf("Operation %s has unexpected Imm %b", test.name, fields.Imm)
			}
		})
	}
}
