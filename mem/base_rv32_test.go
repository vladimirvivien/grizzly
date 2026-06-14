//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package mem

import (
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

func TestBaseMemory_ReadWrite(t *testing.T) {
	size := 1024 * 100
	mem := NewBase(uint64(size))

	// initialize mem
	for i := 0; i < size; i += datapath.XWordBytes {
		if i > mem.GetSize()-datapath.XWordBytes {
			break
		}
		value := datapath.XWord(i * 0x11223344)
		mem.Write(datapath.XWord(i), value, store.Sw.F3)
	}

	// test mem
	for i := 0; i < size; i += datapath.XWordBytes {
		if i > mem.GetSize()-datapath.XWordBytes {
			break
		}
		expected := datapath.XWord(i * 0x11223344)
		val := mem.Read(datapath.XWord(i), load.Lw.F3)
		if val != expected {
			t.Errorf("unexpected value mem[%d]=%d", i, val)
		}
	}
}
