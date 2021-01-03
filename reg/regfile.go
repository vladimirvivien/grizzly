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
	log.Println("regfile: initializing...")
	rs1Pin := r.GetPin(In.RS1Addr)
	rs2Pin := r.GetPin(In.RS2Addr)
	rdPin := r.GetPin(In.RDAddr)
	werfPin := r.GetPin(In.Werf)
	dataPin := r.GetPin(In.Data)

	// Regfile Read Loop
	go func() {
		defer func() {
			close(r.rs1DataOut)
			close(r.rs2DataOut)
		}()

		for {
			// Collect address lines
			addrs := datapath.Collect(rs1Pin, rs2Pin)
			rs1Addr, rs2Addr := addrs[0], addrs[1]

			// send values out
			datapath.Send(
				datapath.Packet{r.read(rs1Addr), r.rs1DataOut},
				datapath.Packet{r.read(rs2Addr), r.rs2DataOut},
			)
		}
	}()

	// Regfile Write Loop
	go func() {
		for {
			// receive write-enble reg file (werf) control
			// then only store value if is set to 1
			select {
			case werf := <-werfPin:
				results := datapath.Collect(rdPin, dataPin)
				if werf == 0 {
					continue
				}
				rdAddr, data := results[0], results[1]
				r.write(rdAddr, data)
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
	log.Printf("regfile: probing addr[%05b]", addr)
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
