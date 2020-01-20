package decoder

import (
	"fmt"
	"testing"

	"github.com/vladimirvivien/grizzly/isa"
)

// TestDecoder_DecodeRFormat tests R-format instructions
func TestDecoder_R(t *testing.T) {
	tests := []struct {
		name       string
		encode     func() isa.Inst
		assess     func(*isa.Fields) error
		shouldFail bool
	}{
		{
			name: isa.Add.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_000_00101_0110011
			},
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
				if fields.Opcode != isa.Add.Opcode || functs != isa.Add.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.Add.Name, fields)
				}
				return nil
			},
		},
		{
			name: isa.Sub.Name,
			encode: func() isa.Inst {
				return 0b0100000_00010_00001_000_00101_0110011
			},
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
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
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
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
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
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
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
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
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
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
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
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
			assess: func(fields *isa.Fields) error {
				functs := isa.EncFuncts(fields.Funct7, fields.Funct3)
				if fields.Opcode != isa.And.Opcode || functs != isa.And.Functs {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", isa.And.Name, fields)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields, err := decode(test.encode())
			if err != nil {
				if !test.shouldFail {
					t.Fatalf("Unexpected error: %s", err)
				}
				t.Logf("Error: %s", err)
			}

			err = test.assess(fields)

			if err != nil {
				if !test.shouldFail {
					t.Fatalf("Unexpected error: %s", err)
				}
				t.Logf("Error: %s", err)
			}
		})
	}
}
