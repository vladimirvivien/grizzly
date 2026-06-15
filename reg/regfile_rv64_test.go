//go:build rv64 || rv64i

package reg

import (
	"math"
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

func TestRegisterFile_Probe_RV64(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8]datapath.XWord
	}{
		{
			name: "probe-64",
			data: map[uint8]datapath.XWord{
				2:  1 << 2,
				4:  1 << 4,
				8:  1 << 8,
				16: 1 << 16,
				31: 1 << 62, // 64-bit value
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := New()
			for addr, val := range test.data {
				reg.file[addr] = val
			}

			for addr, expected := range test.data {
				if val := reg.Probe(addr); val != expected {
					t.Errorf("unexpected val %d probed from addr %d, expected %d", val, addr, expected)
				}
			}
		})
	}
}

func TestRegisterFile_Sideload_RV64(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8]datapath.XWord
	}{
		{
			name: "side-64",
			data: map[uint8]datapath.XWord{
				2:  1 << 2,
				4:  1 << 4,
				8:  1 << 8,
				16: 1 << 16,
				31: 1 << 62, // 64-bit value
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := New()
			for addr, val := range test.data {
				reg.Sideload(addr, val)
			}

			for addr, expected := range test.data {
				if reg.file[addr] != expected {
					t.Errorf("unexpected val %d probed from addr %d, expected %d", reg.file[addr], addr, expected)
				}
			}
		})
	}
}

func TestRegisterFile_Run_FieldsInput_RV64(t *testing.T) {
	reg := New()
	reg.Sideload(1, 0x123456789ABCDEF0)
	reg.Sideload(2, 0x0F0F0F0F0F0F0F0F)
	ch := make(chan []byte)
	reg.Connect(Labels.InFields, ch)
	reg.Connect(Labels.InAluData, make(chan []byte))
	reg.Connect(Labels.InMemData, make(chan []byte))

	go func() {
		ch <- datapath.EncodeOpFields(datapath.OpFields{Opcode: isa.Opcodes.R, Rd: 0b00101, Funct3: 0, Rs1: 0b00001, Rs2: 0b00010, Funct7: 0})
		reg.writeSig <- writeSignal{}
		close(ch)
	}()

	if err := reg.Run(); err != nil {
		t.Fatal(err)
	}

	waiter := make(chan struct{})
	go func() {
		defer close(waiter)
		for output := range reg.GetPin(Labels.OutAluOps) {
			op := datapath.DecodeOp(output)
			if op.Opcode != isa.Opcodes.R {
				t.Errorf("unexpected opcode: %d", op.Opcode)
			}
			if op.AluOperand1 != 0x123456789ABCDEF0 {
				t.Errorf("unexpected AluOperand1: %x", op.AluOperand1)
			}
			if op.AluOperand2 != 0x0F0F0F0F0F0F0F0F {
				t.Errorf("unexpected AluOperand2: %x", op.AluOperand2)
			}
		}
	}()

	select {
	case <-waiter:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("timed out")
	}
}

func TestRegisterFile_Run_AluDataInput_RV64(t *testing.T) {
	tests := []struct {
		name string
		data map[uint8][]byte
	}{
		{
			name: "values",
			data: map[uint8][]byte{
				0:  datapath.EncodeRegData(datapath.RegisterData{Rd: 0, Value: 0}),
				4:  datapath.EncodeRegData(datapath.RegisterData{Rd: 4, Value: math.MaxUint64}),
				8:  datapath.EncodeRegData(datapath.RegisterData{Rd: 8, Value: 0x123456789ABCDEF0}),
				31: datapath.EncodeRegData(datapath.RegisterData{Rd: 31, Value: 0x0F0F0F0F0F0F0F0F}),
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
				<-reg.writeSig
			}
		}()

		if err := reg.Run(); err != nil {
			t.Fatal(err)
		}

		select {
		case <-waiter:
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timed out")
		}

		for addr, stream := range test.data {
			data := datapath.DecodeRegData(stream)
			actual := reg.Probe(addr)
			if addr == 0 {
				if actual != 0 {
					t.Errorf("register x0 must always be 0, got %d", actual)
				}
			} else {
				if actual != data.Value {
					t.Errorf("register x%d: expected %x, got %x", addr, data.Value, actual)
				}
			}
		}
	}
}

func FuzzRegisterFile(f *testing.F) {
	f.Add(uint8(5), uint64(0x123456789ABCDEF0))
	f.Add(uint8(0), uint64(0x9ABCDEF012345678))
	f.Fuzz(func(t *testing.T, regAddr uint8, val uint64) {
		regAddr = regAddr % 32

		reg := New()
		ch := make(chan []byte)
		reg.Connect(Labels.InFields, make(chan []byte))
		reg.Connect(Labels.InMemData, make(chan []byte))
		reg.Connect(Labels.InAluData, ch)

		if err := reg.Run(); err != nil {
			t.Fatal(err)
		}

		go func() {
			ch <- datapath.EncodeRegData(datapath.RegisterData{Rd: regAddr, Value: val})
			close(ch)
		}()

		select {
		case <-reg.writeSig:
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timed out waiting for write signal")
		}

		actual := reg.Probe(regAddr)
		if regAddr == 0 {
			if actual != 0 {
				t.Errorf("x0 must always be 0, got %d", actual)
			}
		} else {
			if actual != val {
				t.Errorf("register x%d: expected %x, got %x", regAddr, val, actual)
			}
		}
	})
}
