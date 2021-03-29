package datapath

import (
	"sync"
)

// Bytestream is a stream of instruction bytes
type Bytestream = <-chan []byte
// Pin is a label for component connection
type Pin string

// Component represents a component on the datapath.
// A component has input/output connection pins identified
// with a label
type Component interface{
	// Connect connects a bytestream to the component pin
	Connect(Pin,Bytestream)
	// GetPin returns bytestream connected to the labeled pin
	GetPin(Pin) Bytestream
	// Starts component
	Run() error
}

// BaseComponent is a base type component
// on which to build others
type BaseComponent struct {
	sync.RWMutex
	pins map[Pin]Bytestream
}

func NewBase() *BaseComponent {
	return &BaseComponent{pins: make(map[Pin]Bytestream)}
}

// GetPin returns a bytestream connected to the labeled pin
func (c *BaseComponent) GetPin(label Pin) Bytestream {
	c.RLock()
	defer c.RUnlock()
	return c.pins[label]
}

// SetPin connects a bytestream to the labeled component pin
func (c *BaseComponent) Connect(label Pin, stream Bytestream) {
	c.Lock()
	defer c.Unlock()
	c.pins[label] = stream
}