package reg

import (
	"fmt"
	"log"

	"github.com/vladimirvivien/grizzly/datapath"
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
	rs1DataOut datapath.Wires
	rs2DataOut datapath.Wires
}

func New() device.Type {
	return newRegister()
}

func newRegister() *RegisterFile {
	r := &RegisterFile{
		file:       make([]datapath.Word, 32, 32),
		rs1DataOut: datapath.MakeWires(),
		rs2DataOut: datapath.MakeWires(),
		Base:       device.NewBase(),
	}

	// wire pin port
	r.SetPin(Out.RS1Data, r.rs1DataOut)
	r.SetPin(Out.RS2Data, r.rs2DataOut)

	return r
}

func (r *RegisterFile) Run() error {
	log.Println("regfile: starting...")
	// rs1
	rs1Addr := r.GetPin(In.RS1Addr)
	rs2Addr := r.GetPin(In.RS2Addr)
	go func() {
		defer close(r.rs1DataOut)
		for {
			select {
			case addr := <-rs1Addr:
				r.rs1DataOut <- r.read(addr)
			case addr := <-rs2Addr:
				r.rs2DataOut <- r.read(addr)
			default:
				continue
			}
		}
	}()

	// write operation blocks until it is received
	dataPin := r.GetPin(In.Data)
	dataAddrPin := r.GetPin(In.RDAddr)
	werfPin := r.GetPin(In.Werf)

	// write loop sequence:
	// 1. writeEnale
	// 2. RD Data
	// 3. RD Addr
	go func() {
		for {
			select {
			case <-werfPin:
				addr := <-dataAddrPin
				data := <-dataPin
				r.write(addr, data)
			default:
				continue
			}
		}
	}()

	return nil
}

func (r *RegisterFile) read(addr uint32) (data uint32) {
	r.RLock()
	defer r.RUnlock()
	if addr == 0 {
		data = 0
	} else {
		data = r.file[addr]
	}
	log.Printf("regfile: read file[%05b] = %032b", addr, data)
	return
}

func (r *RegisterFile) write(addr uint32, data uint32) {
	r.Lock()
	defer r.Unlock()
	if addr == 0 {
		return
	}
	r.file[addr] = data
	log.Printf("regfile: wrote file[%05b]=%032b", addr, data)
}

// Probe is a TEST-ONLY method that is used to read
// values from register address directly.
func (r *RegisterFile) Probe(addr uint32) uint32 {
	return r.read(addr)
}

// SideLoad is TEST-ONLY method used to load values directly into reg
func (r *RegisterFile) SideLoad(addr uint32, val uint32) {
	log.Printf("regfile: sideload addr[%05b]=%032b", addr, val)
	r.write(addr, val)
}

func (r *RegisterFile) Print() {
	fmt.Println()
	for i, v := range r.file {
		fmt.Printf("reg[%d]=%d;", i, v)
	}
	fmt.Println()
}
