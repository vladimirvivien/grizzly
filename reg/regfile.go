package reg

import (
	"sync"

	"github.com/vladimirvivien/grizzly/device"
)

var (
	Wires = struct {
		RS1AddrIn,
		RS2AddrIn,
		RDAddrIn,
		DataIn,
		RS1DataOut,
		RS2DataOut string
	}{
		RS1AddrIn:  "regfile.rs1Addr.in",
		RS2AddrIn:  "regfile.rs2Addr.in",
		RDAddrIn:   "regfile.rdAddr.in",
		DataIn:     "regfile.data.in",
		RS1DataOut: "regfile.rs1Data.out",
		RS2DataOut: "regfile.rs2Data.out",
	}
)

type RegisterFile struct {
	file       []uint32
	rs1AddrIn  device.WiresIn
	rs2AddrIn  device.WiresIn
	rdAddrIn   device.WiresIn
	dataIn     device.WiresIn
	rs1DataOut device.Wires
	rs2DataOut device.Wires
	loadRdy    chan struct{}
	storeRdy   chan struct{}
	enabled    bool
	sync.RWMutex
}

func New() device.Type {
	return newRegister()
}

func newRegister() *RegisterFile {
	return &RegisterFile{
		file:       make([]uint32, 32, 32),
		rs1DataOut: device.MakeWires(),
		rs2DataOut: device.MakeWires(),
		loadRdy:    make(chan struct{}),
		storeRdy:   make(chan struct{}),
	}
}

func (r *RegisterFile) Run() error {
	r.setEnabled(true)

	// process RS1, RS2
	go func() {
		defer func() {
			close(r.rs1DataOut)
			close(r.rs2DataOut)
		}()

		for {
			if !r.IsEnabled() {
				return
			}

			select {
			// wait for store ready signal before writing data value
			case <-r.storeRdy:
				select {
				// write data
				case rdAddr := <-r.rdAddrIn:
					// wait for data line
					select {
					case data := <-r.dataIn:
						r.store(rdAddr, data)
					}
				}

			// wait for load signal before reading rs1 or rs2
			case <-r.loadRdy:
				// read file[rs1]
				go func() {
					addr := <-r.rs1AddrIn
					r.rs1DataOut <- r.load(addr)
				}()

				// read file[rs2]
				go func() {
					addr := <-r.rs2AddrIn
					r.rs2DataOut <- r.load(addr)
				}()
			}
		}
	}()

	return nil
}

func (r *RegisterFile) Port() device.Port {
	return device.Port{
		Wires.RS1AddrIn:  r.rs1AddrIn,
		Wires.RS2AddrIn:  r.rs2AddrIn,
		Wires.RDAddrIn:   r.rdAddrIn,
		Wires.DataIn:     r.dataIn,
		Wires.RS1DataOut: r.rs1DataOut,
		Wires.RS2DataOut: r.rs2DataOut,
	}
}

func (r *RegisterFile) RS1AddrIn(rs1 device.WiresIn) {
	r.rs1AddrIn = rs1
}

func (r *RegisterFile) RS2AddrIn(rs2 device.WiresIn) {
	r.rs2AddrIn = rs2
}

func (r *RegisterFile) RDAddrIn(rd device.WiresIn) {
	r.rdAddrIn = rd
}

func (r *RegisterFile) DataIn(data device.WiresIn) {
	r.dataIn = data
}

func (r *RegisterFile) RS1DataOut() device.WiresOut {
	return r.rs1DataOut
}

func (r *RegisterFile) RS2DataOut() device.WiresOut {
	return r.rs2DataOut
}

func (r *RegisterFile) Load() {
	go func() {
		r.loadRdy <- struct{}{}
	}()
}

func (r *RegisterFile) Store() {
	go func() {
		r.storeRdy <- struct{}{}
	}()
}

// IsEnabled reflects readiness of device.
func (r *RegisterFile) IsEnabled() bool {
	r.RLock()
	defer r.RUnlock()
	return r.enabled
}

func (r *RegisterFile) setEnabled(f bool) {
	r.Lock()
	defer r.Unlock()
	r.enabled = f
}

func (r *RegisterFile) load(addr uint32) uint32 {
	r.RLock()
	defer r.RUnlock()
	return r.file[addr]
}

func (r *RegisterFile) store(addr uint32, data uint32) {
	r.Lock()
	defer r.Unlock()
	r.file[addr] = data
}
