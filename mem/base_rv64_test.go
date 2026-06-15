//go:build rv64 || rv64i

package mem

import (
	"testing"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

func TestBaseMemory_ReadWrite_RV64(t *testing.T) {
	size := 1024
	mem := NewBase(uint64(size))
	mem.SetStore(make([]byte, size))

	// Test writing and reading 64-bit double words (LD/SD)
	mem.Write(0, 0x123456789ABCDEF0, 8) // default/LD size is 8
	val := mem.Read(0, 8)
	if val != 0x123456789ABCDEF0 {
		t.Errorf("expected 0x123456789ABCDEF0, got %x", val)
	}

	// Test 32-bit sign extension (LW)
	mem.Write(8, 0xFFFFFFFF80000000, store.Sw.F3)
	val = mem.Read(8, load.Lw.F3)
	if val != 0xFFFFFFFF80000000 {
		t.Errorf("expected sign extended 0xFFFFFFFF80000000, got %x", val)
	}

	// Test 16-bit unsigned / signed load (LH / LHU)
	mem.Write(16, 0x8000, store.Sh.F3)
	val = mem.Read(16, load.Lh.F3) // signed
	if val != 0xFFFFFFFFFFFF8000 {
		t.Errorf("expected signed extended halfword 0xFFFFFFFFFFFF8000, got %x", val)
	}
	val = mem.Read(16, load.Lhu.F3) // unsigned
	if val != 0x8000 {
		t.Errorf("expected unsigned halfword 0x8000, got %x", val)
	}

	// Test 8-bit unsigned / signed load (LB / LBU)
	mem.Write(24, 0x80, store.Sb.F3)
	val = mem.Read(24, load.Lb.F3) // signed
	if val != 0xFFFFFFFFFFFFFF80 {
		t.Errorf("expected signed extended byte 0xFFFFFFFFFFFFFF80, got %x", val)
	}
	val = mem.Read(24, load.Lbu.F3) // unsigned
	if val != 0x80 {
		t.Errorf("expected unsigned byte 0x80, got %x", val)
	}
}

func TestBaseMemory_AlignmentPanic_RV64(t *testing.T) {
	mem := NewBase(64)
	mem.SetStore(make([]byte, 64))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for misaligned address")
		}
	}()

	mem.Write(1, 42, store.Sb.F3) // misaligned
}

func TestBaseMemory_OutOfBoundsPanic_RV64(t *testing.T) {
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
	f.Add(uint32(8), uint64(0x123456789ABCDEF0), uint8(store.Sb.F3), uint8(load.Lb.F3))
	f.Add(uint32(16), uint64(0xFFFFFFFFFFFFFF80), uint8(store.Sh.F3), uint8(load.Lhu.F3))
	f.Fuzz(func(t *testing.T, offset uint32, val uint64, stOp uint8, ldOp uint8) {
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

		var expected uint64
		switch stOp {
		case store.Sb.F3:
			expected = val & 0xFF
		case store.Sh.F3:
			expected = val & 0xFFFF
		case store.Sw.F3:
			expected = val & 0xFFFFFFFF
		}

		// Fuzz load operations
		switch ldOp % 5 {
		case 0:
			ldOp = load.Lb.F3
			// check signed byte extension
			expectedSign := int64(int8(expected & 0xFF))
			actual := mem.Read(addr, ldOp)
			if int64(actual) != expectedSign {
				t.Errorf("LB: got %x, expected %x", actual, expectedSign)
			}
		case 1:
			ldOp = load.Lbu.F3
			actual := mem.Read(addr, ldOp)
			if uint64(actual) != (expected & 0xFF) {
				t.Errorf("LBU: got %x, expected %x", actual, expected&0xFF)
			}
		case 2:
			ldOp = load.Lh.F3
			expectedSign := int64(int16(expected & 0xFFFF))
			actual := mem.Read(addr, ldOp)
			if int64(actual) != expectedSign {
				t.Errorf("LH: got %x, expected %x", actual, expectedSign)
			}
		case 3:
			ldOp = load.Lhu.F3
			actual := mem.Read(addr, ldOp)
			if uint64(actual) != (expected & 0xFFFF) {
				t.Errorf("LHU: got %x, expected %x", actual, expected&0xFFFF)
			}
		case 4:
			ldOp = load.Lw.F3
			expectedSign := int64(int32(expected & 0xFFFFFFFF))
			actual := mem.Read(addr, ldOp)
			if int64(actual) != expectedSign {
				t.Errorf("LW: got %x, expected %x", actual, expectedSign)
			}
		}
	})
}
