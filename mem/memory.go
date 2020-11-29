package mem

import (
	"encoding/binary"
	"log"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
)

var (
	In = struct {
		Address,
		DataWrite,
		WriteEnable,
		ReadEnable device.PinLabel
	}{
		Address:     "memory.address.in",
		DataWrite:   "memory.datawrite.in",
		WriteEnable: "memory.writeenable.in",
		ReadEnable:  "memory.readenable.in",
	}

	Out = struct {
		DataRead device.PinLabel
	}{
		DataRead: "memory.dataread.out",
	}
)

type Memory struct {
	*device.Base
	store       []byte
	dataReadOut datapath.Wires
}

func New(byteSize uint32) device.Type {
	return newMem(byteSize)
}

func newMem(size uint32) *Memory {
	mem := &Memory{
		Base:        device.NewBase(),
		store:       make([]byte, int(size)),
		dataReadOut: datapath.MakeWires(),
	}
	mem.SetPin(Out.DataRead, mem.dataReadOut)
	return mem
}

func (m *Memory) Run() error {
	log.Println("memory: initializing...")
	addrPin := m.GetPin(In.Address)
	dataWritePin := m.GetPin(In.DataWrite)
	writeEnPin := m.GetPin(In.WriteEnable)
	readEnPin := m.GetPin(In.ReadEnable)

	go func() {
		defer func(){
			close(m.dataReadOut)
		}()

		for {
			addr := <-addrPin
			select {
			case <-readEnPin:
				value := m.read(addr)
				m.dataReadOut <- value
			case <- writeEnPin:
				value := <-dataWritePin
				m.write(addr, value)
			}
		}
	}()

	return nil
}

func (m *Memory) read(addr datapath.Word) (data datapath.Word) {
	m.RLock()
	defer m.RUnlock()

	m.assertAddress(addr)

	switch datapath.Xlen{
	case 32:
		m.assertAlign32(addr)
		buf := m.store[addr:(addr+datapath.XlenBytes)+1]
		data = binary.LittleEndian.Uint32(buf)
		log.Printf("mem: read memory[%032b]=%032b", addr, data)
	default:
	}
	return data
}

func (m *Memory) write(addr, value datapath.Word) {
	m.Lock()
	defer m.Unlock()

	m.assertAddress(addr)

	switch datapath.Xlen {
	case 32:
		m.assertAlign32(addr)
		buf := m.store[addr:(addr+datapath.XlenBytes)+1]
		binary.LittleEndian.PutUint32(buf, value)
		log.Printf("mem: write memory[%032b]=%032b", addr, value)
	case 64:
		if addr & 0x7 > 0 { // 8-byte alignment
			panic("address misaligned")
		}
		binary.LittleEndian.PutUint32(m.store[addr:datapath.XlenBytes], value)
	default:
	}
}

func (m *Memory) assertAddress (addr datapath.Word) {
	if addr > datapath.Word(len(m.store)-4) {
		panic("mem: address out of bound")
	}
	bound := addr+datapath.XlenBytes
	if datapath.Word(len(m.store)) <= bound {
		panic("mem: address out of bound")
	}
}

func (m *Memory) assertAlign32(addr datapath.Word) {
	if addr & 0x3 > 0 { // 4-byte alignment
		panic("mem: address misaligned")
	}
}

func (m *Memory)TestSideLoad(addr datapath.Word, val datapath.Word) {
	m.write(addr, val)
}

func (m *Memory)TestProbe(addr datapath.Word) datapath.Word {
	return m.read(addr)
}