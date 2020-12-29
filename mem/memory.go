package mem

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
)

var (
	In = struct {
		Address,
		Operation,
		DataWrite,
		WriteEnable,
		ReadEnable device.PinLabel
	}{
		Address:     "memory.address.in",
		Operation:   "memory.operation.in",
		DataWrite:   "memory.datawrite.in",
		WriteEnable: "memory.writeenable.in",
		ReadEnable:  "memory.readenable.in",
	}

	Out = struct {
		DataRead device.PinLabel
	}{
		DataRead: "memory.dataread.out",
	}

	// mem operation from instruction Funct3
	Ops = struct {
		Lb,
		Lbu,
		Lh,
		Lhu,
		Lw uint32
	}{
		Lb:  0b_000,
		Lbu: 0b_100,
		Lh:  0b_001,
		Lhu: 0b_101,
		Lw:  0b_110,
	}
)

type Memory struct {
	*device.Base
	state       datapath.Word
	store       []byte
	dataReadOut datapath.Wires
}

func New(byteSize uint64) device.Type {
	return newMem(byteSize)
}

func newMem(size uint64) *Memory {
	mem := &Memory{
		Base:        device.NewBase(),
		store:       make([]byte, size),
		dataReadOut: datapath.MakeWires(),
	}
	mem.SetPin(Out.DataRead, mem.dataReadOut)
	return mem
}

func (m *Memory) Run() error {
	log.Println("memory: initializing...")
	addrPin := m.GetPin(In.Address)
	memOpPin := m.GetPin(In.Operation)
	dataWritePin := m.GetPin(In.DataWrite)
	writeEnPin := m.GetPin(In.WriteEnable)
	readEnPin := m.GetPin(In.ReadEnable)

	go func() {
		defer func() {
			close(m.dataReadOut)
		}()

		for {
			data := datapath.Collect(addrPin, memOpPin)
			addr, memOp := data[0], data[1]

			select {
			case en := <-readEnPin:
				if en != 1 {
					m.dataReadOut <- m.refresh()
					continue
				}
				value := m.read(addr, memOp)
				m.dataReadOut <- value
			case en := <-writeEnPin:
				if en != 1 {
					return
				}
				value := <-dataWritePin
				m.write(addr, value)
			}
		}
	}()

	return nil
}

func (m *Memory) read(addr datapath.Word, op uint32) (data datapath.Word) {
	m.RLock()
	defer m.RUnlock()

	m.assertAddress(addr)

	switch datapath.Xlen {
	case 32:
		m.assertAlign32(addr)
		buf := m.store[addr : (addr+datapath.XlenBytes)+1]
		data = binary.LittleEndian.Uint32(buf)
		log.Printf("mem: read memory[%032b]=%032b", addr, data)
	default:
		panic("mem: unsupported word size")
	}

	// apply operation
	var result datapath.Word

	switch op {
	case Ops.Lb:
		result = datapath.Word(int32(data & 0xFF))
	case Ops.Lbu:
		result = data & 0xFF
	case Ops.Lh:
		result = datapath.Word(int32(data & 0xFFFF))
	case Ops.Lhu:
		result = data & 0xFFFF
	case Ops.Lw:
		result = data
	}

	m.state = result
	return result
}

func (m *Memory) refresh() datapath.Word {
	m.RLock()
	defer m.RUnlock()
	return m.state
}

func (m *Memory) write(addr, value datapath.Word) {
	m.Lock()
	defer m.Unlock()

	m.assertAddress(addr)

	switch datapath.Xlen {
	case 32:
		m.assertAlign32(addr)
		buf := m.store[addr : (addr+datapath.XlenBytes)+1]
		binary.LittleEndian.PutUint32(buf, value)
		log.Printf("mem: write memory[%032b]=%032b", addr, value)
	case 64:
		if addr&0x7 > 0 { // 8-byte alignment
			panic("mem: address misaligned")
		}
		binary.LittleEndian.PutUint32(m.store[addr:datapath.XlenBytes], value)
	default:
	}
}

func (m *Memory) assertAddress(addr datapath.Word) {
	if addr > datapath.Word(len(m.store)-4) {
		panic(fmt.Sprintf("mem: address %032b out of bound", addr))
	}
	bound := addr + datapath.XlenBytes
	if datapath.Word(len(m.store)) <= bound {
		panic(fmt.Sprintf("mem: address %032b out of bound", addr))
	}
}

func (m *Memory) assertAlign32(addr datapath.Word) {
	if addr&0x3 > 0 { // 4-byte alignment
		panic("mem: address misaligned")
	}
}

// TestSideLoad is TEST-ONLY method used to load values directly into memory
func (m *Memory) TestSideLoad(addr datapath.Word, val datapath.Word) {
	m.write(addr, val)
}

// TestProbe is TEST-ONLY method used to read values directly from mem
func (m *Memory) TestProbe(addr datapath.Word) datapath.Word {
	return m.read(addr, Ops.Lw)
}
