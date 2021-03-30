package mem

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

func TestMemory_ReadWrite(t *testing.T) {
	size := 1024 * 100
	mem := New(uint64(size))

	// initialize mem
	for i := 0; i < size; i += datapath.XWordBytes {
		if i > len(mem.store)-datapath.XWordBytes {
			break
		}
		value := datapath.XWord(i * 0x11223344)
		mem.TestSideLoad(datapath.XWord(i), value)
	}

	// test mem
	for i := 0; i < size; i += datapath.XWordBytes {
		if i > len(mem.store)-datapath.XWordBytes {
			break
		}
		expected := datapath.XWord(i * 0x11223344)
		val := mem.TestProbe(datapath.XWord(i))
		if val != expected {
			t.Errorf("unexpected value mem[%d]=%d", i, val)
		}
	}
}

func TestMemory_Run_Read(t *testing.T) {
	size := 1024 * 100
	opCh := make(chan []byte)
	mem := New(uint64(size))
	mem.Connect(Labels.InOperation, opCh)

	if err := mem.Run(); err != nil {
		t.Fatal(err)
	}

	// initialize mem
	for i := 0; i < size; i += datapath.XWordBytes {
		if i > len(mem.store)-datapath.XWordBytes {
			break
		}
		value := datapath.XWord(i * 0x11223344)
		mem.TestSideLoad(datapath.XWord(i), value)
	}

	// send load memory operation
	go func() {
		for i := 0; i < size; i += datapath.XWordBytes {
			if i > len(mem.store)-datapath.XWordBytes {
				break
			}
			opCh <- datapath.EncodeMemOp(datapath.MemOp{
				Opcode: isa.Opcodes.L,
				Rd:     5,
				Funct3: load.Lw.F3,
				Addr:   datapath.XWord(i),
			})
		}
		close(opCh)
	}()

	waiter := make(chan struct{})
	go func() {
		defer close(waiter)
		output := mem.GetPin(Labels.OutRegStore)
		i := 0
		for {
			stream, opened := <-output
			if !opened {
				return
			}
			rs := datapath.DecodeRegStore(stream)
			if rs.Rd != 5 {
				t.Errorf("unexpected regStore.Rd: %d", rs.Rd)
			}
			if rs.Data != datapath.XWord(i*0x11223344) {
				t.Errorf("unexpected data loaded")
			}
			i += datapath.XWordBytes
		}
	}()

	select {
	case <-waiter:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Register operations took too long to complete")
	}
}

func TestMemory_Run_Write(t *testing.T) {
	size := 1024 * 100
	opCh := make(chan []byte)
	mem := New(uint64(size))
	mem.Connect(Labels.InOperation, opCh)

	if err := mem.Run(); err != nil {
		t.Fatal(err)
	}

	waiter := make(chan struct{})
	// store
	go func() {
		for i := 0; i < size; i += datapath.XWordBytes {
			if i > len(mem.store)-datapath.XWordBytes {
				break
			}
			opCh <- datapath.EncodeMemOp(datapath.MemOp{
				Opcode: isa.Opcodes.S,
				Funct3: store.Sw.F3,
				Addr:   datapath.XWord(i),
				Data:   datapath.XWord(i * 0x11223344),
			})
		}
		close(opCh)
		close(waiter)
	}()

	select {
	case <-waiter:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Register operations took too long to complete")
	}

	// test mem
	for i := 0; i < size; i += datapath.XWordBytes {
		if i > len(mem.store)-datapath.XWordBytes {
			break
		}
		expected := datapath.XWord(i * 0x11223344)
		val := mem.TestProbe(datapath.XWord(i))
		if val != expected {
			t.Errorf("unexpected value mem[%d]=%d", i, val)
		}
	}
}
