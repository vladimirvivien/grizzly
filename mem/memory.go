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
		Lw,
		Sb,
		Sh,
		Sw uint32
	}{
		Lb:  0b_000,
		Lbu: 0b_100,
		Lh:  0b_001,
		Lhu: 0b_101,
		Lw:  0b_110,
		Sb:  0b_000,
		Sh:  0b_001,
		Sw:  0b_010,
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

		rcvr := datapath.NewReceiver("mem")
		for {
			data := rcvr.R(addrPin, memOpPin)
			addr, memOp := data[0], data[1]

			select {
			case en := <-readEnPin:
				// This case is hit for both non-memory and Load instructions.
				// When non-mem instructions are executed, read-enable = 0 and
				// it allows data to stream through the memory and back
				// to register to avoid blocking.
				if en == 0 {
					m.dataReadOut <- m.refresh()
					log.Println("mem: read enabled=false")
					continue
				}

				// For Load mem instructions, read-enable = 1, then the memory
				// is actually carried out.
				value := m.read(addr, memOp)
				m.dataReadOut <- value
				log.Printf("mem:read enable=true; addr=%032b; data=%032b; op=%03b", addr, value, memOp)

			case <-writeEnPin:
				// Write enabled is sent only when instructions is for store
				value := <-dataWritePin
				m.write(addr, value, memOp)

				// Always send data to the mem-output to ensure bit stream continuity
				// out of the mem component to avoid deadlock
				m.dataReadOut <- m.refresh()
				log.Printf("mem:write addr=%032b; data=%032b; op=%03b", addr, value, memOp)
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
		buf := m.store[addr : addr+datapath.XlenBytes]
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

func (m *Memory) write(addr, value datapath.Word, op uint32) {
	m.Lock()
	defer m.Unlock()

	m.assertAddress(addr)

	// apply store operation
	var data datapath.Word
	switch op {
	case Ops.Sb:
		data = value & 0xFF
	case Ops.Sh:
		data = value & 0xFFFF
	default:
		data = value
	}

	switch datapath.Xlen {
	case 32:
		m.assertAlign32(addr)
		buf := m.store[addr : addr+datapath.XlenBytes]
		binary.LittleEndian.PutUint32(buf, data)
		log.Printf("mem: write memory[%032b]=%032b", addr, data)
	case 64:
		if addr&0x7 > 0 { // 8-byte alignment
			panic("mem: address misaligned")
		}
		binary.LittleEndian.PutUint32(m.store[addr:datapath.XlenBytes], data)
		log.Printf("mem: write memory[%064b]=%064b", addr, data)
	default:
	}
}

func (m *Memory) assertAddress(addr datapath.Word) {
	if addr > datapath.Word(len(m.store)-4) {
		panic(fmt.Sprintf("mem: address %032b out of bound", addr))
	}
	bound := addr + datapath.XlenBytes
	if bound > datapath.Word(len(m.store)) {
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
	m.write(addr, val, Ops.Lw)
}

// TestProbe is TEST-ONLY method used to read values directly from mem
func (m *Memory) TestProbe(addr datapath.Word) datapath.Word {
	return m.read(addr, Ops.Lw)
}
