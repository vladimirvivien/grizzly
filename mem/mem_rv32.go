package mem

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa/load"
	"github.com/vladimirvivien/grizzly/isa/store"
)

type Memory struct {
	sync.RWMutex
	store []byte
}

func New(size uint64) *Memory {
	return &Memory{store: make([]byte, size, size)}
}

func (m *Memory) read(addr datapath.XWord, f3 uint8) datapath.XWord {
	m.RLock()
	defer m.RUnlock()

	m.assertAddress(addr)
	m.assertAlignment(addr)
	//buf := m.store[addr : addr+datapath.XWordBytes]
	data := binary.LittleEndian.Uint32(m.store[addr:])

	// apply operation
	var result datapath.XWord

	switch f3 {
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

func (m *Memory) write(addr, value datapath.XWord, f3 uint8) {
	m.Lock()
	defer m.Unlock()

	m.assertAddress(addr)
	m.assertAlignment(addr)

	// apply store operation
	var data datapath.XWord
	switch f3 {
	case store.Sb.F3:
		data = value & 0xFF
	case store.Sh.F3:
		data = value & 0xFFFF
	default:
		data = value
	}

	//buf := m.store[addr : addr+datapath.XWordBytes]
	binary.LittleEndian.PutUint32(m.store[addr:], data)
}

// assert address boundaries
func (m *Memory) assertAddress(addr datapath.XWord) {
	if addr > datapath.XWord(len(m.store)-datapath.XWordBytes) {
		panic(fmt.Sprintf("mem: address %032b out of bound", addr))
	}
	bound := addr + datapath.XWordLen
	if bound > datapath.XWord(len(m.store)) {
		panic(fmt.Sprintf("mem: address %032b out of bound", addr))
	}
}

func (m *Memory) assertAlignment(addr datapath.XWord) {
	if addr&0x3 > 0 { // 4-byte alignment
		panic("mem: address misaligned")
	}

	//if addr&0x7 > 0 { // 8-byte alignment
	//	panic("mem: address misaligned")
	//}
}

// TestSideLoad is TEST-ONLY method used to load values directly into memory
func (m *Memory) TestSideLoad(addr datapath.XWord, val datapath.XWord) {
	m.write(addr, val, store.Sw.F3)
}

// TestProbe is TEST-ONLY method used to read values directly from mem
func (m *Memory) TestProbe(addr datapath.XWord) datapath.XWord {
	return m.read(addr, load.Lw.F3)
}

