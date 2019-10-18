package device

type Wires = chan uint32

func MakeWires() Wires {
	return make(chan uint32)
}

type WiresIn = <-chan uint32
type WiresOut = WiresIn
type Datapath []Wires

type Port map[string]<-chan uint32

type Type interface {
	Run() error
	Port() Port
}
