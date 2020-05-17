package device

// TODO
// investigate using []byte to represent instruction stream between components
// Then use encoding/binary package to narrow []byte values to numeric instruction
// see https://play.golang.org/p/py_Uv9zSXWv
//
// This would allow Grizzly to support multi-size instruction (32, 64, 128, etc) and
// also support for compressed instructions at runtime.

type Datapath = []byte

type Wires = chan uint32

func MakeWires() Wires {
	return make(chan uint32)
}

type WiresIn = <-chan uint32
type WiresOut = WiresIn

type Pin = <-chan uint32
type PinLabel = string
type Pins = map[PinLabel]Pin

type Type interface {
	Run() error
	GetPins() Pins
	GetPin(PinLabel) Pin
	SetPin(PinLabel, Pin)
}
