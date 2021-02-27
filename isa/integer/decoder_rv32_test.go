package integer

import (
	"fmt"
	"testing"

	"github.com/vladimirvivien/grizzly/isa"
)

// TestDecoder tests R-format integer instructions
//
// fn7     RD2   RD1   fn3 RS    OPCODE
// 0000000_00010_00001_000_00101_0110011
func TestDecoder(t *testing.T) {
	tests := []struct {
		name       string
		encode     func() isa.Inst
		assess     func(*Fields) error
		shouldFail bool
	}{
		{
			name: Add.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_000_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != Add.Opcode {
					return fmt.Errorf("Unexpected Opcode %b in %s: %#v", fields.Opcode, Add.Name, fields)
				}
				if functs != Add.Functs {
					return fmt.Errorf("Unexpected Operation for %s: %#v", Add.Name, fields)
				}
				if fields.Rd != 0b00101 {
					return fmt.Errorf("isa.Instruction %s has wrong RD field value: %#v", Add.Name, fields)
				}
				if fields.Rs1 != 0b00001 {
					return fmt.Errorf("isa.Instruction %s has wrong RS1 field value: %#v", Add.Name, fields)
				}
				if fields.Rs2 != 0b00010 {
					return fmt.Errorf("isa.Instruction %s has wrong RS2 field value: %#v", Add.Name, fields)
				}
				return nil
			},
		},
		{
			name: Sub.Name,
			encode: func() isa.Inst {
				return 0b0100000_00010_00001_000_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != Sub.Opcode || functs != Sub.Functs {
					return fmt.Errorf("isa.Instruction %s has wrong field values: %#v", Sub.Name, fields)
				}
				return nil
			},
		},
		{
			name: Sll.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_001_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != Sll.Opcode || functs != Sll.Functs {
					return fmt.Errorf("isa.Instruction %s has wrong field values: %#v", Sll.Name, fields)
				}
				return nil
			},
		},
		{
			name: Slt.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_010_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != Slt.Opcode || functs != Slt.Functs {
					return fmt.Errorf("isa.Instruction %s has wrong field values: %#v", Slt.Name, fields)
				}
				return nil
			},
		},
		{
			name: Srl.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_101_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != Srl.Opcode || functs != Srl.Functs {
					return fmt.Errorf("isa.Instruction %s has wrong field values: %#v", Srl.Name, fields)
				}
				return nil
			},
		},
		{
			name: Sra.Name,
			encode: func() isa.Inst {
				return 0b0100000_00010_00001_101_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != Sra.Opcode || functs != Sra.Functs {
					return fmt.Errorf("isa.Instruction %s has wrong field values: %#v", Sra.Name, fields)
				}
				return nil
			},
		},
		{
			name: Or.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_110_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != Or.Opcode || functs != Or.Functs {
					return fmt.Errorf("isa.Instruction %s has wrong field values: %#v", Or.Name, fields)
				}
				return nil
			},
		},
		{
			name: And.Name,
			encode: func() isa.Inst {
				return 0b0000000_00010_00001_111_00101_0110011
			},
			assess: func(fields *Fields) error {
				functs := fields.Functs()
				if fields.Opcode != And.Opcode || functs != And.Functs {
					return fmt.Errorf("isa.Instruction %s has wrong field values: %#v", And.Name, fields)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Decode(test.encode())
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
