package mem

import (
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestMemory_ReadWrite(t *testing.T) {
	size := 1024 * 100
	mem := New(uint64(size))

	// initialize mem
	for i := 0; i < size; i += datapath.XWordLen {
		if size <= i+datapath.XWordLen {
			break
		}
		value := datapath.XWord(i*4)
		mem.TestSideLoad(datapath.XWord(i), value)
	}

	// test mem
	for i := 0; i < size; i += datapath.XWordLen {
		if size <= i+datapath.XWordLen {
			break
		}
		expected := datapath.XWord(i*4)
		val := mem.TestProbe(datapath.XWord(i))
		if val != expected {
			t.Errorf("unexpected value mem[%032b]=%032b", i, val)
		}
	}
}