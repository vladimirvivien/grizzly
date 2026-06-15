//go:build rv64 || rv64i

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
			fields: datapath.OpFields{Imm: 0b1000100010010010, Rd: 0b00001},
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

func FuzzDecodeJump(f *testing.F) {
	f.Add(uint32(0b00001001001100001000_00001_1101111))
	f.Add(uint32(0b000010010011_00101_000_00001_1100111))
	f.Fuzz(func(t *testing.T, instVal uint32) {
		inst := datapath.XWord(instVal)
		opcode := inst & 0x7F
		if opcode != 0x6F && opcode != 0x67 {
			return
		}

		fields := Decode(inst)

		if fields.Opcode != uint8(opcode) {
			t.Errorf("opcode mismatch: got %x, expected %x", fields.Opcode, opcode)
		}
		if fields.Rd > 31 || fields.Rs1 > 31 {
			t.Errorf("registers must be <= 31: rd=%d, rs1=%d", fields.Rd, fields.Rs1)
		}

		if opcode == 0x6F {
			imm20 := (inst >> 31) & 0x1
			imm10_1 := (inst >> 21) & 0x3FF
			imm11 := (inst >> 20) & 0x1
			imm19_12 := (inst >> 12) & 0xFF
			val := imm10_1 | imm11 << 10 | imm19_12 << 11 | imm20 << 19
			offset := val << 1
			if (offset & 0x100000) != 0 {
				offset |= 0xffe00000
			}
			if fields.Imm != uint32(offset) {
				t.Errorf("JAL imm mismatch: got %x, expected %x", fields.Imm, offset)
			}
		} else if opcode == 0x67 {
			imm := (inst >> 20) & 0xFFF
			if (imm & 0x800) != 0 {
				imm |= 0xfffff000
			}
			if fields.Imm != uint32(imm) {
				t.Errorf("JALR imm mismatch: got %x, expected %x", fields.Imm, imm)
			}
		}
	})
}
