package device

// Base is a base type on which to build other device
type Base struct {
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
	return c.pins[label]
}

// SetPin sets specified labeled pin to a pin
func (c *Base) SetPin(label PinLabel, pin Pin) {
	c.pins[label] = pin
}
