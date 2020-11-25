package device

import (
	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
)

type Pin = datapath.WireRcvd
type PinLabel = string
type Pins = map[PinLabel]Pin

type Type interface {
	Run() error
	GetPins() Pins
	GetPin(PinLabel) Pin
	SetPin(PinLabel, Pin)
}

type ClockedType interface {
	Type
	SetClock (clock.Clock)
}