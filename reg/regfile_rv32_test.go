package reg

import (
	"math"
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
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

// test instruction operation fields input
func TestRegisterFile_Run_FieldsInput(t *testing.T) {
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
				ch := make(chan []byte)
				reg.Connect(Labels.InFields, ch)
				go func() {
					// 0b0000000_00010_00001_000_00101_0110011
					ch <- datapath.EncodeOpFields(datapath.OpFields{Opcode: isa.Opcodes.R, Rd: 0b00101, Funct3: 0, Rs1: 0b00001, Rs2: 0b00010, Funct7: 0})
					reg.writeSig <- writeSignal{}
					close(ch)
				}()
				return reg
			},
			eval: func(t *testing.T, reg *RegisterFile) {
				for output := range reg.GetPin(Labels.OutAluOps) {
					op := datapath.DecodeOp(output)
					if op.Opcode != isa.Opcodes.R {
						t.Errorf("unexpected ALUParam.opcode: %d", op.Opcode)
					}

					if op.Rd != 0b00101 {
						t.Errorf("unexpected ALUParam.Rd: %d", op.Rd)
					}

					if op.AluOp != alu.Ops.Add {
						t.Errorf("unexpected aluOp: %d", op.AluOp)
					}

					if op.AluOperand1 != math.MaxInt32 {
						t.Errorf("unexpected ALUParam.AluOperand1: %d", op.AluOperand1)
					}
					if op.AluOperand2 != math.MaxInt8 {
						t.Errorf("unexpected ALUParam.AluOperand2: %d", op.AluOperand2)
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
				ch := make(chan []byte)
				reg.Connect(Labels.InFields, ch)
				go func() {
					// 0b0000000_00010_00001_000_00101_0110011 (add)
					ch <- datapath.EncodeOpFields(datapath.OpFields{Opcode: 0b0110011, Rd: 0b00101, Funct3: 0, Rs1: 0b00001, Rs2: 0b00010, Funct7: 0})
					reg.writeSig <- writeSignal{}
					// 0b000000000010_00001_000_00101_0010011 (addi)
					ch <- datapath.EncodeOpFields(datapath.OpFields{Imm: 0b000000000010, Rs2: 0b00010, Rs1: 0b00001, Funct3: 0b000, Rd: 0b00101, Opcode: 0b0010011})
					reg.writeSig <- writeSignal{}
					// 0b0100000_11011_01101_101_00101_0010011 (srai)
					ch <- datapath.EncodeOpFields(datapath.OpFields{Shift: 0b11011, Funct7: 0b0100000, Rs1: 0b01101, Funct3: 0b101, Rd: 0b00101, Opcode: 0b0010011})
					reg.writeSig <- writeSignal{}
					close(ch)
				}()
				return reg
			},
			eval: func(t *testing.T, reg *RegisterFile) {
				// instruction 1
				op := datapath.DecodeOp(<-reg.GetPin(Labels.OutAluOps))
				if op.Opcode != isa.Opcodes.R {
					t.Errorf("unexpected ALUParam.opcode: %d", op.Opcode)
				}

				// instruction 2
				op = datapath.DecodeOp(<-reg.GetPin(Labels.OutAluOps))
				if op.Opcode != 0b0010011 {
					t.Errorf("unexpected opcode: %d", op.Opcode)
				}
				if op.AluOperand1 != math.MaxInt32 {
					t.Errorf("unexpected AluOperand1: %d", op.AluOperand1)
				}
				if op.AluOperand2 != 2 {
					t.Errorf("unexpected AluOperand2: %d", op.AluOperand2)
				}
				if op.AluOp != alu.Ops.Add {
					t.Errorf("unexpected aluOp: %d", op.AluOp)
				}

				// instruction 3
				op = datapath.DecodeOp(<-reg.GetPin(Labels.OutAluOps))
				if op.Opcode != 0b0010011 {
					t.Errorf("unexpected ALUParam.opcode: %d", op.Opcode)
				}
				if op.AluOperand1 != 32 {
					t.Errorf("unexpected ALUParam.AluOperand1: %d", op.AluOperand1)
				}
				if op.AluOperand2 != 0b11011 {
					t.Errorf("unexpected ALUParam.AluOperand2: %d", op.AluOperand2)
				}
				if op.AluOp != alu.Ops.Sra {
					t.Errorf("unexpected aluOp: %d", op.AluOp)
				}

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := test.setup(t)
			reg.Connect(Labels.InAluData, make(chan []byte))
			reg.Connect(Labels.InMemData, make(chan []byte))

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

// Test register data input from alu
func TestRegisterFile_Run_AluDataInput(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8][]byte
	}{
		{
			name: "values",
			data: map[uint8][]byte{
				0:  datapath.EncodeRegData(datapath.RegisterData{Rd: 0, Value: 0}),
				4:  datapath.EncodeRegData(datapath.RegisterData{Rd: 4, Value: math.MaxUint32}),
				8:  datapath.EncodeRegData(datapath.RegisterData{Rd: 8, Value: math.MaxInt8}),
				16: datapath.EncodeRegData(datapath.RegisterData{Rd: 16, Value: math.MaxInt16}),
				31: datapath.EncodeRegData(datapath.RegisterData{Rd: 31, Value: math.MaxInt32}),
			},
		},
	}

	for _, test := range tests {
		reg := New()
		ch := make(chan []byte)
		reg.Connect(Labels.InFields, make(chan []byte))
		reg.Connect(Labels.InMemData, make(chan []byte))
		reg.Connect(Labels.InAluData, ch)

		waiter := make(chan struct{})
		go func() {
			defer close(ch)
			defer close(waiter)
			for _, data := range test.data {
				ch <- data
				<-reg.writeSig // unblock signal
			}
		}()

		if err := reg.Run(); err != nil {
			t.Fatal(err)
		}

		select {
		case <-waiter:
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Register operations took too long to complete")
		}

		// assess
		for addr, stream := range test.data {
			data := datapath.DecodeRegData(stream)
			actual := reg.Probe(addr)
			if actual != data.Value {
				t.Errorf("expecting RegisterData.Vallue %d, got %d", data.Value, actual)
			}
		}
	}
}


// Test register data input from memory
func TestRegisterFile_Run_MemDataInput(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8][]byte
	}{
		{
			name: "values",
			data: map[uint8][]byte{
				0:  datapath.EncodeRegData(datapath.RegisterData{Rd: 0, Value: 0}),
				4:  datapath.EncodeRegData(datapath.RegisterData{Rd: 4, Value: math.MaxUint32}),
				8:  datapath.EncodeRegData(datapath.RegisterData{Rd: 8, Value: math.MaxInt8}),
				16: datapath.EncodeRegData(datapath.RegisterData{Rd: 16, Value: math.MaxInt16}),
				31: datapath.EncodeRegData(datapath.RegisterData{Rd: 31, Value: math.MaxInt32}),
			},
		},
	}

	for _, test := range tests {
		reg := New()
		ch := make(chan []byte)
		reg.Connect(Labels.InFields, make(chan []byte))
		reg.Connect(Labels.InAluData, make(chan []byte))
		reg.Connect(Labels.InMemData, ch)

		waiter := make(chan struct{})
		go func() {
			defer close(ch)
			defer close(waiter)
			for _, data := range test.data {
				ch <- data
				<-reg.writeSig // unblock signal
			}
		}()

		if err := reg.Run(); err != nil {
			t.Fatal(err)
		}

		select {
		case <-waiter:
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Register operations took too long to complete")
		}

		// assess
		for addr, stream := range test.data {
			data := datapath.DecodeRegData(stream)
			actual := reg.Probe(addr)
			if actual != data.Value {
				t.Errorf("expecting RegisterData.Vallue %d, got %d", data.Value, actual)
			}
		}
	}
}