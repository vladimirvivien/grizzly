package mem

import (
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/device"
)

var (
	In = struct {
		Address,
		WriteData,
		WriteEnable,
		ReadEnable device.PinLabel
	}{
		Address:     "memory.address.in",
		WriteData:   "memory.writedata.in",
		WriteEnable: "memory.writeenable.in",
		ReadEnable:  "memory.readenable.in",
	}

	Out = struct {
		ReadData device.PinLabel
	}{
		ReadData: "memory.readdata.out",
	}
)

type Memory struct {
	*device.Base
	store       []byte
	dataOut datapath.Wires
}