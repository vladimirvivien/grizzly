package decoder

import (
	"fmt"
	"testing"

	"github.com/vladimirvivien/grizzly/inst"
)

// TestDecoder_DecodeRFormat tests R-format instructions
func TestDecoder_R(t *testing.T) {
	tests := []struct {
		name       string
		encode     func() inst.Type
		assess     func(*inst.Fields) error
		shouldFail bool
	}{
		{
			name: inst.Add.Name,
			encode: func() inst.Type {
				return 0b0000000_00010_00001_000_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.Add.Opcode ||
					fields.Funct7 != inst.Add.Funct7 ||
					fields.Funct3 != inst.Add.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.Add.Name, fields)
				}
				return nil
			},
		},
		{
			name: inst.Sub.Name,
			encode: func() inst.Type {
				return 0b0100000_00010_00001_000_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.Sub.Opcode ||
					fields.Funct7 != inst.Sub.Funct7 ||
					fields.Funct3 != inst.Sub.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.Sub.Name, fields)
				}
				return nil
			},
		},
		{
			name: inst.Sll.Name,
			encode: func() inst.Type {
				return 0b0000000_00010_00001_001_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.Sll.Opcode ||
					fields.Funct7 != inst.Sll.Funct7 ||
					fields.Funct3 != inst.Sll.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.Sll.Name, fields)
				}
				return nil
			},
		},
		{
			name: inst.Slt.Name,
			encode: func() inst.Type {
				return 0b0000000_00010_00001_010_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.Slt.Opcode ||
					fields.Funct7 != inst.Slt.Funct7 ||
					fields.Funct3 != inst.Slt.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.Slt.Name, fields)
				}
				return nil
			},
		},
		{
			name: inst.Srl.Name,
			encode: func() inst.Type {
				return 0b0000000_00010_00001_101_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.Srl.Opcode ||
					fields.Funct7 != inst.Srl.Funct7 ||
					fields.Funct3 != inst.Srl.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.Srl.Name, fields)
				}
				return nil
			},
		},
		{
			name: inst.Sra.Name,
			encode: func() inst.Type {
				return 0b0100000_00010_00001_101_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.Sra.Opcode ||
					fields.Funct7 != inst.Sra.Funct7 ||
					fields.Funct3 != inst.Sra.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.Sra.Name, fields)
				}
				return nil
			},
		},
		{
			name: inst.Or.Name,
			encode: func() inst.Type {
				return 0b0000000_00010_00001_110_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.Or.Opcode ||
					fields.Funct7 != inst.Or.Funct7 ||
					fields.Funct3 != inst.Or.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.Or.Name, fields)
				}
				return nil
			},
		},
		{
			name: inst.And.Name,
			encode: func() inst.Type {
				return 0b0000000_00010_00001_111_00101_0110011
			},
			assess: func(fields *inst.Fields) error {
				if fields.Opcode != inst.And.Opcode ||
					fields.Funct7 != inst.And.Funct7 ||
					fields.Funct3 != inst.And.Funct3 {
					return fmt.Errorf("Instruction %s has wrong field values: %#v", inst.And.Name, fields)
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
