package jump

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
			name:   Jal.Name,
			inst:   0b00001001001100001000_00001_1101111,
			fields: datapath.OpFields{Imm: 0b00001001001100001000, Rd: 0b00001},
		},
		{
			name:   Jalr.Name,
			inst:   0b000010010011_00101_000_00001_1100111,
			fields: datapath.OpFields{Imm: 0b000010010011, Rs1: 0b00101, Rd: 0b00001},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.inst)
			switch fields.Opcode {
			case isa.Opcodes.J, isa.Opcodes.JI:
			default:
				t.Errorf("Unexpected Opcode %05b for op %s: %#v", fields.Opcode, test.name, fields)
			}
			if fields.Rd != test.fields.Rd {
				t.Errorf("Operation %s has unexpected Rd %b", test.name, fields.Rd)
			}
			if fields.Funct3 != test.fields.Funct3 {
				t.Errorf("Operation %s has unexpected Func3 %b", test.name, fields.Funct3)
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
