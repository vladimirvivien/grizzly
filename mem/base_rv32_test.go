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

func TestBaseMemory_AlignmentPanic(t *testing.T) {
	mem := NewBase(64)
	mem.SetStore(make([]byte, 64))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for misaligned address")
		}
	}()

	mem.Write(1, 42, store.Sb.F3) // misaligned
}

func TestBaseMemory_OutOfBoundsPanic(t *testing.T) {
	mem := NewBase(16)
	mem.SetStore(make([]byte, 16))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for out of bounds address")
		}
	}()

	mem.Write(16, 42, store.Sb.F3) // out of bounds
}

func FuzzMemory(f *testing.F) {
	f.Add(uint32(8), uint32(0x12345678), uint8(store.Sb.F3), uint8(load.Lb.F3))
	f.Add(uint32(16), uint32(0xFFFF9ABC), uint8(store.Sh.F3), uint8(load.Lhu.F3))
	f.Fuzz(func(t *testing.T, offset uint32, val uint32, stOp uint8, ldOp uint8) {
		memSize := uint64(128)
		mem := NewBase(memSize)
		mem.SetStore(make([]byte, memSize))

		// Enforce 4-byte alignment
		addr := datapath.XWord(offset % 96) &^ 3

		var writeSize datapath.XWord
		switch stOp % 3 {
		case 0:
			stOp = store.Sb.F3
			writeSize = 1
		case 1:
			stOp = store.Sh.F3
			writeSize = 2
		default:
			stOp = store.Sw.F3
			writeSize = 4
		}

		// Prevent out of bounds
		if addr+writeSize > datapath.XWord(memSize) {
			return
		}

		mem.Write(addr, datapath.XWord(val), stOp)

		var expected uint32
		switch stOp {
		case store.Sb.F3:
			expected = val & 0xFF
		case store.Sh.F3:
			expected = val & 0xFFFF
		case store.Sw.F3:
			expected = val
		}

		// Fuzz load operations
		switch ldOp % 5 {
		case 0:
			ldOp = load.Lb.F3
			// check signed byte extension
			expectedSign := int32(int8(expected & 0xFF))
			actual := mem.Read(addr, ldOp)
			if int32(actual) != expectedSign {
				t.Errorf("LB: got %x, expected %x", actual, expectedSign)
			}
		case 1:
			ldOp = load.Lbu.F3
			actual := mem.Read(addr, ldOp)
			if uint32(actual) != (expected & 0xFF) {
				t.Errorf("LBU: got %x, expected %x", actual, expected&0xFF)
			}
		case 2:
			ldOp = load.Lh.F3
			expectedSign := int32(int16(expected & 0xFFFF))
			actual := mem.Read(addr, ldOp)
			if int32(actual) != expectedSign {
				t.Errorf("LH: got %x, expected %x", actual, expectedSign)
			}
		case 3:
			ldOp = load.Lhu.F3
			actual := mem.Read(addr, ldOp)
			if uint32(actual) != (expected & 0xFFFF) {
				t.Errorf("LHU: got %x, expected %x", actual, expected&0xFFFF)
			}
		case 4:
			ldOp = load.Lw.F3
			actual := mem.Read(addr, ldOp)
			if uint32(actual) != expected {
				t.Errorf("LW: got %x, expected %x", actual, expected)
			}
		}
	})
}

