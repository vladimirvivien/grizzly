//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

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
		size: size,
	}
}

func (m *BaseMemory)SetStore(store []byte) {
	m.store = store
}

func (m *BaseMemory) GetSize() int {
	return len(m.store)
}

func (m *BaseMemory) Read(addr datapath.XWord, memOp uint8) datapath.XWord {
	m.RLock()
	defer m.RUnlock()

	m.AssertAddress(addr)
	m.AssertAlignment(addr)
	data := binary.LittleEndian.Uint32(m.store[addr:])

	// apply operation
	var result datapath.XWord

	switch memOp {
	case load.Lb.F3:
		result = datapath.XWord(int32(data & 0xFF))
	case load.Lbu.F3:
		result = data & 0xFF
	case load.Lh.F3:
		result = datapath.XWord(int32(data & 0xFFFF))
	case load.Lhu.F3:
		result = data & 0xFFFF
	case load.Lw.F3:
		result = data
	}

	return result
}

func (m *BaseMemory) Write(addr, value datapath.XWord, memOp uint8) {
	m.Lock()
	defer m.Unlock()

	m.AssertAddress(addr)
	m.AssertAlignment(addr)

	// apply store operation
	var data datapath.XWord
	switch memOp {
	case store.Sb.F3:
		data = value & 0xFF
	case store.Sh.F3:
		data = value & 0xFFFF
	default:
		data = value
	}

	binary.LittleEndian.PutUint32(m.store[addr:], data)
}

// assert address boundaries
func (m *BaseMemory) AssertAddress(addr datapath.XWord) {
	if addr > datapath.XWord(len(m.store)-datapath.XWordBytes) {
		panic(fmt.Sprintf("mem: address %d out of bound", addr))
	}
}

func (m *BaseMemory) AssertAlignment(addr datapath.XWord) {
	if addr&0x3 > 0 { // 4-byte alignment
		panic("mem: address misaligned")
	}
	//if addr&0x7 > 0 { // 8-byte alignment
	//	panic("mem: address misaligned")
	//}
}
