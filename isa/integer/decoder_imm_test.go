package integer

import (
	"testing"

	"github.com/vladimirvivien/grizzly/isa"
)

// TestDecoderImm tests integer immediate instructions (I-format)
//
// Immediate ops
// immOut[32:20]   RS1   fn3 RD    OPCODE
// 000000000010_00001_000_00101_0110011
//
// Shift ops
// fn7     Shift RS1   fn3 RD    OPCODE
// 0000000_00010_00001_000_00101_0110011
func TestDecodeImm(t *testing.T) {
	tests := []struct {
		name                 string
		inst                 isa.Inst
		shift, imm, funct7   uint32
		rs2, rs1, funct3, rd uint32
		shouldFail           bool
	}{
		{
			name:   Addi.Name,
			inst:   0b000000000010_00001_000_00101_0010011,
			imm:    0b000000000010,
			rs1:    0b00001,
			funct3: 0b000,
			rd:     0b00101,
		},
		{
			name:   Slli.Name,
			inst:   0b0000000_00010_00001_001_00101_0010011,
			funct7: 0b0000000,
			shift:  0b00010,
			rs1:    0b00001,
			funct3: 0b001,
			rd:     0b00101,
		},
		{
			name:   Slti.Name,
			inst:   0b000000000010_00001_010_10101_0010011,
			imm:    0b000000000010,
			rs1:    0b00001,
			funct3: 0b010,
			rd:     0b10101,
		},
		{
			name:   Srli.Name,
			inst:   0b0000000_01010_00001_101_01101_0010011,
			funct7: 0b0000000,
			shift:  0b01010,
			rs1:    0b00001,
			funct3: 0b101,
			rd:     0b01101,
		},
		{
			name:   Srai.Name,
			inst:   0b0100000_11011_01101_101_00101_0010011,
			funct7: 0b0100000,
			shift:  0b11011,
			rs1:    0b01101,
			funct3: 0b101,
			rd:     0b00101,
		},
		{
			name:   Ori.Name,
			inst:   0b010001000010_00001_110_00101_0010011,
			imm:    0b010001000010,
			rs1:    0b00001,
			funct3: 0b110,
			rd:     0b00101,
		},
		{
			name:   Andi.Name,
			inst:   0b010000000010_10101_111_10101_0010011,
			imm:    0b010000000010,
			rs1:    0b10101,
			funct3: 0b111,
			rd:     0b10101,
		},
	}

	riopcode := uint32(0b0010011)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := DecodeImm(test.inst)
			if fields.Opcode != riopcode {
				t.Errorf("Unexpected Opcode %b for op %s: %#v", fields.Opcode, test.name, fields)
			}
			if fields.Funct3 != test.funct3 {
				t.Errorf("Operation %s has unexpected Funct3 %b", test.name, fields.Functs())
			}
			if fields.Rd != test.rd {
				t.Errorf("Operation %s has unexpected RD %b", test.name, fields.Rd)
			}
			if fields.Rs1 != test.rs1 {
				t.Errorf("Operation %s has unexpected RS1 %b", test.name, fields.Rs1)
			}
			if fields.Shift != test.shift {
				t.Errorf("Operation %s has unexpected shift amt %b", test.name, fields.Imm)
			}
			if fields.Imm != test.imm {
				t.Errorf("Operation %s has unexpected Imm %b", test.name, fields.Imm)
			}
		})
	}
}
