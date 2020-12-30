package store

import (
	"testing"

	"github.com/vladimirvivien/grizzly/isa"
)

// TestDecode tests store instructions of form
// 31.......25.....20.....15...12.....07.......0
//   I[11:5]  RS2    RS1    fn3  I[4:0] OPCODE
//   0000000__00000__00000__xxx__00000__0000011
//
func TestDecode(t *testing.T) {
	tests := []struct {
		name             string
		inst             isa.Inst
		hiImm            uint32
		rs2, rs1, funct3 uint32
		loImm            uint32
		imm              uint32
	}{
		{
			name:   Sb.Name,
			inst:   0b0000100_10011_00001_000_00001_0100011,
			hiImm:  0b0000100,
			rs2:    0b10011,
			rs1:    0b00001,
			funct3: 0b000,
			loImm:  0b00001,
			imm:    0b000010000001,
		},
		{
			name:   Sh.Name,
			inst:   0b0010100_11011_01001_001_00101_0100011,
			hiImm:  0b0010100,
			rs2:    0b11011,
			rs1:    0b01001,
			funct3: 0b001,
			loImm:  0b00101,
			imm:    0b001010000101,
		},
		{
			name:   Sw.Name,
			inst:   0b0010100_11011_01001_010_00101_0100011,
			hiImm:  0b0010100,
			rs2:    0b11011,
			rs1:    0b01001,
			funct3: 0b010,
			loImm:  0b00101,
			imm:    0b001010000101,
		},
	}

	sopcode := uint32(0b0100011)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.inst)
			if fields.Opcode != sopcode {
				t.Errorf("Unexpected Opcode %05b for op %s: %#v", fields.Opcode, test.name, fields)
			}
			if fields.Funct3 != test.funct3 {
				t.Errorf("Operation %s has unexpected Funct3 %b", test.name, fields.Funct3)
			}
			if fields.loImm != test.loImm {
				t.Errorf("Operation %s has unexpected RD %b", test.name, fields.Rd)
			}
			if fields.Rs1 != test.rs1 {
				t.Errorf("Operation %s has unexpected RS1 %b", test.name, fields.Rs1)
			}
			if fields.Rs2 != test.rs2 {
				t.Errorf("Operation %s has unexpected RS2 %b", test.name, fields.Rs2)
			}
			if fields.Imm != test.imm {
				t.Errorf("Operation %s has unexpected Imm %b", test.name, fields.Imm)
			}
		})
	}
}
