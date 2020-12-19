package mem

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestMemory_ReadWrite(t *testing.T) {
	size := 1024
	mem := newMem(uint64(size))

	for i := 0; i < size; i += datapath.XlenBytes {
		if size <= i+datapath.XlenBytes {
			break
		}
		mem.TestSideLoad(datapath.Word(i), datapath.Word(i*4))
		val := mem.TestProbe(datapath.Word(i))
		if val != datapath.Word(i*4) {
			t.Errorf("unexpected value mem[%032b]=%0b32", i, val)
		}
	}
}

func TestMemory_New(t *testing.T) {
	size := uint64(1024)
	mem := New(size)
	if mem.GetPin(Out.DataRead) == nil {
		t.Error("mem: pin DataRead not set")
	}
	memory := mem.(*Memory)
	if len(memory.store) != 1024 {
		t.Error("mem: memory not initialized")
	}
}

func TestMemory_Run(t *testing.T) {
	size := 1024 * 2
	addr := datapath.MakeWires()
	wen := datapath.MakeWires()
	ren := datapath.MakeWires()
	data := datapath.MakeWires()
	mem := New(uint64(size))
	mem.SetPin(In.Address, addr)
	mem.SetPin(In.DataWrite, data)
	mem.SetPin(In.WriteEnable, wen)
	mem.SetPin(In.ReadEnable, ren)

	if err := mem.Run(); err != nil {
		t.Fatal(err)
	}

	waiter := make(chan struct{})
	go func() {
		defer close(waiter)
		for i := 0; i < size; i += datapath.XlenBytes {
			if size <= i+datapath.XlenBytes {
				break
			}
			addr <- datapath.Word(i)
			wen <- 1
			data <- datapath.Word(i * 4)
		}
	}()

	select {
	case <-waiter:
		dataOut := mem.GetPin(Out.DataRead)
		for i := 0; i < size; i+=datapath.XlenBytes{
			if size <= i+datapath.XlenBytes {
				break
			}
			addr <- datapath.Word(i)
			ren <- 1
			val := <- dataOut
			if val != datapath.Word(i*4) {
				t.Errorf("mem: unexpected value memory[%032b]=%032b", i, val)
			}
		}
	case <-time.After(15*time.Millisecond):
		t.Fatal("mem: took too long to initialize")
	}
}
