package load

import (
	"testing"

	"github.com/vladimirvivien/grizzly/isa"
)

//TestDecode load isa instruction
//   imm[31:20]   RS1   fn3 RD    OPCODE
//   000000000000_00000_xxx_00000_0000011
func TestDecode(t *testing.T) {
	tests := []struct {
		name                 string
		inst                 isa.Inst
		imm                  uint32
		rs2, rs1, funct3, rd uint32
	}{
		{
			name:   Lb.Name,
			inst:   0b000000010011_00001_000_00001_0000011,
			imm:    0b000000010011,
			rs1:    0b00001,
			funct3: 0b000,
			rd:     0b00001,
		},
		{
			name:   Lbu.Name,
			inst:   0b100000011011_01001_100_11001_0000011,
			imm:    0b100000011011,
			rs1:    0b01001,
			funct3: 0b100,
			rd:     0b11001,
		},
		{
			name:   Lh.Name,
			inst:   0b000010001011_01011_001_11101_0000011,
			imm:    0b000010001011,
			rs1:    0b01011,
			funct3: 0b001,
			rd:     0b11101,
		},
		{
			name:   Lhu.Name,
			inst:   0b100010001011_01011_101_00101_0000011,
			imm:    0b100010001011,
			rs1:    0b01011,
			funct3: 0b101,
			rd:     0b00101,
		},
		{
			name:   Lw.Name,
			inst:   0b111111001011_11011_110_00101_0000011,
			imm:    0b111111001011,
			rs1:    0b11011,
			funct3: 0b110,
			rd:     0b00101,
		},
	}

	lopcode := uint32(0b0000011)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.inst)
			if fields.Opcode != lopcode {
				t.Errorf("Unexpected Opcode %05b for op %s: %#v", fields.Opcode, test.name, fields)
			}
			if fields.Funct3 != test.funct3 {
				t.Errorf("Operation %s has unexpected Funct3 %b", test.name, fields.Funct3)
			}
			if fields.Rd != test.rd {
				t.Errorf("Operation %s has unexpected RD %b", test.name, fields.Rd)
			}
			if fields.Rs1 != test.rs1 {
				t.Errorf("Operation %s has unexpected RS1 %b", test.name, fields.Rs1)
			}
			if fields.Imm != test.imm {
				t.Errorf("Operation %s has unexpected Imm %b", test.name, fields.Imm)
			}
		})
	}
}
