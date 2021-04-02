package pc

import (
	"fmt"
	"time"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
)

var (
	Labels = struct {
		InPcOp     datapath.Pin
		OutCounter datapath.Pin
	}{
		InPcOp:     datapath.Pin("pc.in.pc_operation"),
		OutCounter: datapath.Pin("pc.out.counter"),
	}
)

// PC represents the program counter component.
// It is clocked and its cycle trigers downstream
// components.
type PC struct {
	*datapath.BaseComponent
	clock    *clock.Clock
	counter  datapath.XWord
	transfer chan datapath.XWord
	out      chan []byte
}

func New() *PC {
	pc := &PC{
		BaseComponent: datapath.NewBase(),
		clock:         clock.New(time.Microsecond),
		transfer:      make(chan datapath.XWord),
		out:           make(chan []byte),
	}
	pc.Connect(Labels.OutCounter, pc.out)
	return pc
}

func (pc *PC) Run() error {
	opCh := pc.GetPin(Labels.InPcOp)
	if opCh == nil {
		return fmt.Errorf("pc: missing input: %s", Labels.InPcOp)
	}

	go func() {
		for stream := range opCh {
			op := datapath.DecodePcOp(stream)
			if op.Jump > 0 {
				pc.counter = op.PC
			} else {
				pc.counter = pc.counter + 4
			}
			pc.transfer <- pc.counter
		}
	}()

	// output loop
	// generates pc
	go func() {
		defer close(pc.out)
		for range pc.clock.Ticks() {
			pc.out <- datapath.EncodePC(<-pc.transfer)
		}
	}()

	return nil
}
