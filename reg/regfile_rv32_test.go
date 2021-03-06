package reg

import (
	"math"
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestRegisterFile_Probe(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8]datapath.XWord
	}{
		{
			name: "probe-32",
			data: map[uint8]datapath.XWord{
				2:  1 << 2,
				4:  1 << 4,
				8:  1 << 8,
				16: 1 << 16,
				31: 1 << 31,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := New()
			for addr, val := range test.data {
				reg.file[addr] = val
			}

			for addr, _ := range test.data {
				if val := reg.Probe(addr); val != 1<<addr {
					t.Errorf("unexpected val %d probed from addr %d", val, addr)
				}
			}
		})
	}
}

func TestRegisterFile_Sideload(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8]datapath.XWord
	}{
		{
			name: "side-32",
			data: map[uint8]datapath.XWord{
				2:  1 << 2,
				4:  1 << 4,
				8:  1 << 8,
				16: 1 << 16,
				31: 1 << 31,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := New()
			for addr, val := range test.data {
				reg.Sideload(addr, val)
			}

			for addr, _ := range test.data {
				if reg.file[addr] != 1<<addr {
					t.Errorf("unexpected val %d probed from addr %d", reg.file[addr], addr)
				}
			}
		})
	}
}

// test instruction operation input
func TestRegisterFile_Run_Op(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*testing.T) *RegisterFile
		eval  func(*testing.T, *RegisterFile)
	}{
		{
			name: "single instructions",
			setup: func(t *testing.T) *RegisterFile {
				reg := New()
				reg.Sideload(1, math.MaxInt32)
				reg.Sideload(2, math.MaxInt8)
				ch := make(chan datapath.OpFields)
				reg.OpInput(ch)
				go func() {
					// 0b0000000_00010_00001_000_00101_0110011
					ch <- datapath.OpFields{Opcode: 0b0110011, Rd: 0b00101, Funct3: 0, Rs1: 0b00001, Rs2: 0b00010, Funct7: 0}
					close(ch)
				}()
				return reg
			},
			eval: func(t *testing.T, reg *RegisterFile) {
				for param := range reg.AluParams() {
					if param.Opcode != 0b0110011 {
						t.Errorf("unexpected ALUParam.opcode: %d", param.Opcode)
					}

					if param.Rd != 0b00101 {
						t.Errorf("unexpected ALUParam.Rd: %d", param.Rd)
					}

					if param.Funct3 != 0 {
						t.Errorf("unexpected ALUParam.Funct3: %d", param.Funct3)
					}

					if param.Funct7 != 0 {
						t.Errorf("unexpected ALUParam.Funct7: %d", param.Funct7)
					}

					if param.Op1 != math.MaxInt32 {
						t.Errorf("unexpected ALUParam.Op1: %d", param.Op1)
					}
					if param.Op2 != math.MaxInt8 {
						t.Errorf("unexpected ALUParam.Op2: %d", param.Op2)
					}
				}
			},
		},

		{
			name: "multiple instructions",
			setup: func(t *testing.T) *RegisterFile {
				reg := New()
				reg.Sideload(1, math.MaxInt32)
				reg.Sideload(2, math.MaxInt8)
				reg.Sideload(13, 32)
				ch := make(chan datapath.OpFields)
				reg.OpInput(ch)
				go func() {
					// 0b0000000_00010_00001_000_00101_0110011 (add)
					ch <- datapath.OpFields{Opcode: 0b0110011, Rd: 0b00101, Funct3: 0, Rs1: 0b00001, Rs2: 0b00010, Funct7: 0}
					// 0b000000000010_00001_000_00101_0010011 (addi)
					ch <- datapath.OpFields{Imm: 0b000000000010, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b000, Rd: 0b00101, Opcode: 0b0010011}
					// 0b0100000_11011_01101_101_00101_0010011 (srai)
					ch <- datapath.OpFields{Shift: 0b11011, Funct7: 0b0100000, Rs1: 0b01101, Funct3: 0b101, Rd: 0b00101, Opcode: 0b0010011}
					close(ch)
				}()
				return reg
			},
			eval: func(t *testing.T, reg *RegisterFile) {
				// instruction 1
				param := <-reg.AluParams()
				if param.Opcode != 0b0110011 {
					t.Errorf("unexpected ALUParam.opcode: %d", param.Opcode)
				}

				// instruction 2
				param = <-reg.AluParams()
				if param.Opcode != 0b0010011 {
					t.Errorf("unexpected ALUParam.opcode: %d", param.Opcode)
				}
				if param.Op1 != math.MaxInt32 {
					t.Errorf("unexpected ALUParam.Op1: %d", param.Op1)
				}
				if param.Op2 != 2 {
					t.Errorf("unexpected ALUParam.Op2: %d", param.Op2)
				}

				// instruction 3
				param = <-reg.AluParams()
				if param.Opcode != 0b0010011 {
					t.Errorf("unexpected ALUParam.opcode: %d", param.Opcode)
				}
				if param.Funct7 != 0b0100000 {
					t.Errorf("unexpected ALUParam.Func7: %d", param.Funct7)
				}
				if param.Op1 != 32 {
					t.Errorf("unexpected ALUParam.Op1: %d", param.Op1)
				}
				if param.Op2 != 0b11011 {
					t.Errorf("unexpected ALUParam.Op2: %d", param.Op2)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := test.setup(t)
			reg.DataInput(make(chan datapath.RegisterData))

			if err := reg.Run(); err != nil {
				t.Fatal(err)
			}
			waiter := make(chan struct{})
			go func() {
				defer close(waiter)
				test.eval(t, reg)
			}()

			select {
			case <-waiter:
			case <-time.After(5 * time.Millisecond):
				t.Fatal("Register operations took too long to complete")
			}
		})
	}
}

// Test register data input
func TestRegisterFile_Run_Data(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8]datapath.RegisterData
	}{
		{
			name: "values",
			data: map[uint8]datapath.RegisterData{
				0: {Rd: 0, Value: 0},
				4:{Rd:4,Value:math.MaxUint32},
				8:{Rd:8,Value:math.MaxInt8},
				16:{Rd:16,Value:math.MaxInt16},
				31:{Rd:31,Value:math.MaxInt32},
			},
		},
	}

	for _, test := range tests {
		reg := New()
		ch := make(chan datapath.RegisterData)
		reg.OpInput(make(chan datapath.OpFields))
		reg.DataInput(ch)

		waiter := make(chan struct{})
		go func() {
			defer close(ch)
			defer close(waiter)
			for _, data := range test.data {
				ch <- data
			}
		}()

		if err := reg.Run(); err != nil {
			t.Fatal(err)
		}

		select {
		case <-waiter:
		case <-time.After(5 * time.Millisecond):
			t.Fatal("Register operations took too long to complete")
		}

		// assess
		for addr, data := range test.data {
			if reg.Probe(addr) != data.Value {
				t.Errorf("unexpected RegisterData.Vallue: %d", data.Value)
			}
		}

	}
}
