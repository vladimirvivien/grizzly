package ctrlunit

import (
	"fmt"
	"testing"

	"github.com/vladimirvivien/grizzly/isa"
)

// TestDecoder_R tests R-format instructions of the form
//
// fn7     RD2   RD1   fn3 RS    OPCODE
// 0000000_00010_00001_000_00101_0110011
func TestDecodeR(t *testing.T) {
	tests := []struct {
		name       string
		encode     func() isa.Inst
		assess     func(*isa.RFields) error
		shouldFail bool
	}{
		{
			name: isa.Add.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_000_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.Add.Opcode {
					return fmt.Errorf("Unexpected Opcode %b in %s: %#v", fields.Opcode, isa.Add.Name, fields)
				}
				if functs != isa.Add.Functs {
					return fmt.Errorf("Unexpected Operation for %s: %#v", isa.Add.Name, fields)
				}
				if fields.Rd != 0b00101 {
					return fmt.Errorf("Instruction %s has wrong RD field value: %#v", isa.Add.Name, fields)
				}
				if fields.Rs1 != 0b00001 {
					return fmt.Errorf("Instruction %s has wrong RS1 field value: %#v", isa.Add.Name, fields)
				}
				if fields.Rs2 != 0b00010 {
					return fmt.Errorf("Instruction %s has wrong RS2 field value: %#v", isa.Add.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.Sub.Name,
			encode: func() isa.Inst {
				return 0b0100000_00010_00001_000_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.Sub.Opcode || functs != isa.Sub.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.Sub.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.Sll.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_001_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.Sll.Opcode || functs != isa.Sll.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.Sll.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.Slt.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_010_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.Slt.Opcode || functs != isa.Slt.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.Slt.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.Srl.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_101_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.Srl.Opcode || functs != isa.Srl.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.Srl.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.Sra.Name,
			encode: func() isa.Inst {
				return 0b0100000_00010_00001_101_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.Sra.Opcode || functs != isa.Sra.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.Sra.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.Or.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_110_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.Or.Opcode || functs != isa.Or.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.Or.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.And.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_111_00101_0110011
			},
			assess: func(fields *isa.RFields) error {
				functs := fields.Functs()
				if fields.Opcode != isa.And.Opcode || functs != isa.And.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.And.Name, fields)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := decodeR(test.encode())
			err := test.assess(fields)
			if err != nil {
				if !test.shouldFail {
					t.Fatalf("Unexpected error: %s", err)
				}
				t.Logf("Error: %s", err)
			}
		})
	}
}
