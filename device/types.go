package device

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
