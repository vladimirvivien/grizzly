package reg

import (
	"sync"

	"github.com/vladimirvivien/grizzly/device"
)

var (
	In = struct {
		RS1Addr,
		RS2Addr,
		RDAddr,
		Data device.PinLabel
	}{
		RS1Addr: "regfile.rs1Addr.in",
		RS2Addr: "regfile.rs2Addr.in",
		RDAddr:  "regfile.rdAddr.in",
		Data:    "regfile.data.in",
	}

	Out = struct {
		RS1Data,
		RS2Data device.PinLabel
	}{
		RS1Data: "regfile.rs1Data.out",
		RS2Data: "regfile.rs2Data.out",
	}
)

type RegisterFile struct {
	file       []uint32
	rs1DataOut device.Wires
	rs2DataOut device.Wires
	pins       device.Pins
	sync.RWMutex
}

func New() device.Type {
	return newRegister()
}

func newRegister() *RegisterFile {
	r := &RegisterFile{
		file:       make([]uint32, 32, 32),
		rs1DataOut: device.MakeWires(),
		rs2DataOut: device.MakeWires(),
		pins:       make(device.Pins),
	}

	// wire pin port
	r.pins = device.Pins{
		Out.RS1Data: r.rs1DataOut,
		Out.RS2Data: r.rs2DataOut,
	}

	return r
}

func (r *RegisterFile) Run() error {
	// rs1
	go func() {
		defer close(r.rs1DataOut)
		for {
			select {
			case addr := <-r.pins[In.RS1Addr]:
				r.rs1DataOut <- r.read(addr)
			}
		}
	}()

	// rs2
	go func() {
		defer close(r.rs2DataOut)
		for {
			select {
			case addr := <-r.pins[In.RS2Addr]:
				r.rs2DataOut <- r.read(addr)
			}
		}
	}()

	// data, rd:
	// for this to work in sequential
	// circuit, data must be specified
	// prior to rd
	go func() {
		for {
			select {
			case data := <-r.pins[In.Data]:
				addr := <-r.pins[In.RDAddr]
				r.write(addr, data)
			}
		}
	}()

	return nil
}

func (r *RegisterFile) GetPins() device.Pins {
	return r.pins
}

func (r *RegisterFile) GetPin(label device.PinLabel) device.Pin {
	return r.pins[label]
}

func (r *RegisterFile) SetPin(label device.PinLabel, pin device.Pin) {
	r.pins[label] = pin
}

func (r *RegisterFile) read(addr uint32) uint32 {
	r.RLock()
	defer r.RUnlock()
	return r.file[addr]
}

func (r *RegisterFile) write(addr uint32, data uint32) {
	r.Lock()
	defer r.Unlock()
	r.file[addr] = data
}
