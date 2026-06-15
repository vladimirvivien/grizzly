//go:build rv64 || rv64i

package branch

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
			name:   Beq.Name,
			inst:   0b0_000100_10011_00001_000_1010_1_1100011,
			fields: datapath.OpFields{Opcode: isa.Opcodes.B, Funct3: Beq.F3, Rs1: 0b00001, Rs2: 0b10011, Imm: 0b0_1_000100_1010},
		},
		{
			name:   Bne.Name,
			inst:   0b1_000100_10011_00001_001_1010_0_1100011,
			fields: datapath.OpFields{Opcode: isa.Opcodes.B, Funct3: Bne.F3, Rs1: 0b00001, Rs2: 0b10011, Imm: 0b1_0_000100_1010},
		},
		{
			name:   Blt.Name,
			inst:   0b1_100100_10011_00001_100_1011_0_1100011,
			fields: datapath.OpFields{Opcode: isa.Opcodes.B, Funct3: Blt.F3, Rs1: 0b00001, Rs2: 0b10011, Imm: 0b1_0_100100_1011},
		},
		{
			name:   Bge.Name,
			inst:   0b1_100100_10011_00001_101_1011_0_1100011,
			fields: datapath.OpFields{Opcode: isa.Opcodes.B, Funct3: Bge.F3, Rs1: 0b00001, Rs2: 0b10011, Imm: 0b1_0_100100_1011},
		},
		{
			name:   Bltu.Name,
			inst:   0b1_111111_10011_00001_110_1011_0_1100011,
			fields: datapath.OpFields{Opcode: isa.Opcodes.B, Funct3: Bltu.F3, Rs1: 0b00001, Rs2: 0b10011, Imm: 0b1_0_111111_1011},
		},
		{
			name:   Bgeu.Name,
			inst:   0b1_111111_10011_00001_111_1011_0_1100011,
			fields: datapath.OpFields{Opcode: isa.Opcodes.B, Funct3: Bgeu.F3, Rs1: 0b00001, Rs2: 0b10011, Imm: 0b1_0_111111_1011},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.inst)
			if fields.Opcode != isa.Opcodes.B {
				t.Errorf("Unexpected Opcode %05b for op %s: %#v", fields.Opcode, test.name, fields)
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
				t.Errorf("Operation %s has unexpected Imm %012b", test.name, fields.Imm)
			}
		})
	}
}

func FuzzDecodeBranch(f *testing.F) {
	// Add seed corpus
	f.Add(uint32(0b0_000100_10011_00001_000_1010_1_1100011))
	f.Fuzz(func(t *testing.T, instVal uint32) {
		inst := datapath.XWord(instVal)
		opcode := inst & 0x7F
		if opcode != 0x63 {
			return
		}

		fields := Decode(inst)

		if fields.Opcode != 0x63 {
			t.Errorf("expected opcode 0x63, got %x", fields.Opcode)
		}
		if fields.Rs1 > 31 || fields.Rs2 > 31 {
			t.Errorf("registers must be <= 31: rs1=%d, rs2=%d", fields.Rs1, fields.Rs2)
		}
		if fields.Funct3 > 7 {
			t.Errorf("funct3 must be <= 7: %d", fields.Funct3)
		}

		// Reconstruction checks
		rs1 := uint8((inst >> 15) & 0x1F)
		rs2 := uint8((inst >> 20) & 0x1F)
		f3 := uint8((inst >> 12) & 0x7)

		imm11 := (inst >> 7) & 0x1
		imm4_1 := (inst >> 8) & 0xF
		imm10_5 := (inst >> 25) & 0x3F
		imm12 := (inst >> 31) & 0x1
		expectedImm := imm4_1 | imm10_5 << 4 | imm11 << 10 | imm12 << 11

		if fields.Rs1 != rs1 {
			t.Errorf("rs1 mismatch: got %d, expected %d", fields.Rs1, rs1)
		}
		if fields.Rs2 != rs2 {
			t.Errorf("rs2 mismatch: got %d, expected %d", fields.Rs2, rs2)
		}
		if fields.Funct3 != f3 {
			t.Errorf("funct3 mismatch: got %d, expected %d", fields.Funct3, f3)
		}
		if fields.Imm != uint32(expectedImm) {
			t.Errorf("imm mismatch: got %d, expected %d", fields.Imm, expectedImm)
		}
	})
}
