package reg

import (
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/device"
)

var (
	In = struct {
		RS1Addr,
		RS2Addr,
		RDAddr,
		Data,
		Werf device.PinLabel
	}{
		RS1Addr: "regfile.rs1Addr.in",
		RS2Addr: "regfile.rs2Addr.in",
		RDAddr:  "regfile.rdAddr.in",
		Data:    "regfile.data.in",
		Werf:    "regfile.werf.in",
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
	*device.Base
	file       []uint32
	rs1DataOut device.Wires
	rs2DataOut device.Wires
	wready     device.Wires
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
		wready:     device.MakeWires(),
		Base:       device.NewBase(),
	}

	// wire pin port
	r.SetPin(Out.RS1Data, r.rs1DataOut)
	r.SetPin(Out.RS2Data, r.rs2DataOut)

	return r
}

func (r *RegisterFile) Run() error {
	// rs1
	go func() {
		defer close(r.rs1DataOut)
		for {
			select {
			case addr := <-r.GetPin(In.RS1Addr):
				r.rs1DataOut <- r.read(addr)
			}
		}
	}()

	// rs2
	go func() {
		defer close(r.rs2DataOut)
		for {
			select {
			case addr := <-r.GetPin(In.RS2Addr):
				r.rs2DataOut <- r.read(addr)
			}
		}
	}()

	// wePin - receives write-enable signal
	// write operation blocks until it is received
	dataPin := r.GetPin(In.Data)
	dataAddrPin := r.GetPin(In.RDAddr)
	werfPin := r.GetPin(In.Werf)

	// write loop sequence:
	// 1. writeEnale
	// 2. RD Data
	// 3. RD Addr
	go func() {
		defer close(r.wready)
		for {
			select {
			case <-werfPin:
				addr := <-dataAddrPin
				data := <-dataPin
				r.write(addr, data)
				// signal write ready
				go func() {
					r.wready <- 1
				}()
			}
		}
	}()

	return nil
}

func (r *RegisterFile) read(addr uint32) uint32 {
	r.RLock()
	defer r.RUnlock()
	if addr == 0 {
		return 0
	}
	return r.file[addr]
}

func (r *RegisterFile) write(addr uint32, data uint32) {
	r.Lock()
	defer r.Unlock()
	if addr == 0 {
		return
	}
	r.file[addr] = data
}

// Probe is a test-only method that blocks until wready
// then reads the specified address
// if wready is triggered by a previous write, this blocks
// indefinitely
func (r *RegisterFile) Probe(addr uint32) uint32 {
	<-r.wready
	return r.read(addr)
}

// SideLoad is test-only method used to load values directly into reg
func (r *RegisterFile) SideLoad(addr uint32, val uint32) {
	r.write(addr, val)
}

func (r *RegisterFile) Print() {
	for i, v := range r.file {
		fmt.Printf("file[%x] = %b\n", i, v)
	}
}
