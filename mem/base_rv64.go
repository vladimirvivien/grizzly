//go:build rv64 || rv64i

package mem

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

type BaseMemory struct {
	sync.RWMutex
	size uint64
	store []byte
	outReg chan []byte
}

func NewBase(size uint64) *BaseMemory {
	return &BaseMemory{
		size:  size,
		store: make([]byte, size),
	}
}

func (m *BaseMemory) SetStore(store []byte) {
	m.store = store
}

func (m *BaseMemory) GetSize() int {
	return len(m.store)
}

func (m *BaseMemory) Read(addr datapath.XWord, memOp uint8) datapath.XWord {
	m.RLock()
	defer m.RUnlock()

	var size datapath.XWord
	switch memOp {
	case load.Lb.F3, load.Lbu.F3:
		size = 1
	case load.Lh.F3, load.Lhu.F3:
		size = 2
	case load.Lw.F3:
		size = 4
	default:
		size = 8
	}

	m.AssertAddress(addr, size)
	m.AssertAlignment(addr)
	
	// Ensure we don't slice out of array bounds before binary read
	if addr + 8 > datapath.XWord(len(m.store)) {
		// If we're reading less than 8 bytes near EOF, we pad the slice to safely decode
		temp := make([]byte, 8)
		copy(temp, m.store[addr:])
		data := binary.LittleEndian.Uint64(temp)
		return m.applyLoadOp(data, memOp)
	}

	data := binary.LittleEndian.Uint64(m.store[addr:])
	return m.applyLoadOp(data, memOp)
}

func (m *BaseMemory) applyLoadOp(data uint64, memOp uint8) datapath.XWord {
	var result datapath.XWord
	switch memOp {
	case load.Lb.F3:
		result = datapath.XWord(int64(int8(data & 0xFF)))
	case load.Lbu.F3:
		result = data & 0xFF
	case load.Lh.F3:
		result = datapath.XWord(int64(int16(data & 0xFFFF)))
	case load.Lhu.F3:
		result = data & 0xFFFF
	case load.Lw.F3:
		result = datapath.XWord(int64(int32(data & 0xFFFFFFFF)))
	default:
		result = data
	}
	return result
}

func (m *BaseMemory) Write(addr, value datapath.XWord, memOp uint8) {
	m.Lock()
	defer m.Unlock()

	var size datapath.XWord
	switch memOp {
	case store.Sb.F3:
		size = 1
	case store.Sh.F3:
		size = 2
	case store.Sw.F3:
		size = 4
	default:
		size = 8
	}

	m.AssertAddress(addr, size)
	m.AssertAlignment(addr)

	if size == 8 {
		binary.LittleEndian.PutUint64(m.store[addr:], value)
	} else {
		// For sub-doubleword writes, we write only the required bytes
		for i := datapath.XWord(0); i < size; i++ {
			m.store[addr+i] = byte((value >> (i * 8)) & 0xFF)
		}
	}
}

func (m *BaseMemory) AssertAddress(addr datapath.XWord, size datapath.XWord) {
	if addr+size > datapath.XWord(len(m.store)) {
		panic(fmt.Sprintf("mem: address %d with size %d out of bound", addr, size))
	}
}

func (m *BaseMemory) AssertAlignment(addr datapath.XWord) {
	if addr&0x3 > 0 { // 4-byte alignment
		panic("mem: address misaligned")
	}
}
