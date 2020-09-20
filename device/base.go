package device

import (
	"sync"
)

// Base is a base type on which to build other device
type Base struct {
	sync.RWMutex
	pins Pins
}

func NewBase() *Base {
	return &Base{
		pins: make(Pins),
	}
}

// GetPins returns a map of device pins
func (c *Base) GetPins() Pins {
	return c.pins
}

// GetPin returns a device pin by label
func (c *Base) GetPin(label PinLabel) Pin {
	c.RLock()
	defer c.RUnlock()
	return c.pins[label]
}

// SetPin sets specified labeled pin to a pin
func (c *Base) SetPin(label PinLabel, pin Pin) {
	c.Lock()
	defer c.Unlock()
	c.pins[label] = pin
}
