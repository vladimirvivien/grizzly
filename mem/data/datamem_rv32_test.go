package data

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)



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
		if i > mem.GetSize()-datapath.XWordBytes {
			break
		}
		value := datapath.XWord(i * 0x11223344)
		mem.TestSideLoad(datapath.XWord(i), value)
	}

	// send load memory operation
	go func() {
		for i := 0; i < size; i += datapath.XWordBytes {
			if i > mem.GetSize()-datapath.XWordBytes {
				break
			}
			opCh <- datapath.EncodeMemOp(datapath.MemOp{
				Opcode: isa.Opcodes.L,
				Rd:     5,
				Op:     load.Lw.F3,
				Addr:   datapath.XWord(i),
			})
		}
		close(opCh)
	}()

	waiter := make(chan struct{})
	go func() {
		defer close(waiter)
		output := mem.GetPin(Labels.OutRegData)
		i := 0
		for {
			stream, opened := <-output
			if !opened {
				return
			}
			rs := datapath.DecodeRegData(stream)
			if rs.Rd != 5 {
				t.Errorf("unexpected regStore.Rd: %d", rs.Rd)
			}
			if rs.Value != datapath.XWord(i*0x11223344) {
				t.Errorf("unexpected data loaded")
			}
			i += datapath.XWordBytes
		}
	}()

	select {
	case <-waiter:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("DataMemory operations took too long to complete")
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
			if i > mem.GetSize()-datapath.XWordBytes {
				break
			}
			opCh <- datapath.EncodeMemOp(datapath.MemOp{
				Opcode: isa.Opcodes.S,
				Op:     store.Sw.F3,
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
		t.Fatal("DataMemory operations took too long to complete")
	}

	// test mem
	for i := 0; i < size; i += datapath.XWordBytes {
		if i > mem.GetSize()-datapath.XWordBytes {
			break
		}
		expected := datapath.XWord(i * 0x11223344)
		val := mem.TestProbe(datapath.XWord(i))
		if val != expected {
			t.Errorf("unexpected value mem[%d]=%d", i, val)
		}
	}
}
