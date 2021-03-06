package device

import (
	"log"

	"github.com/vladimirvivien/grizzly/datapath"
)

// Mux uses selector pin as index to select value from
// a collection of pins (channels)
func Mux(name string, selectorPin Pin, inPins ...Pin) Pin {
	output := datapath.MakeWires()
	go func() {
		defer close(output)
		rcvr := datapath.Receiver{Name:name}
		var words []datapath.XWord
		for {
			select {
			case sel := <-selectorPin:
				log.Printf("mux %s: select-line: %d", name, sel)
				words = rcvr.R(inPins...)
				val := words[sel]
				log.Printf("mux %s: {sel:%d, lines:%d, value:%032b}", name, sel, len(inPins), val)
				output <- val
			}
		}
	}()
	log.Printf("mux %s: 1-to-%d created", name, len(inPins))
	return output
}

// Connector creates size number of output pins and
// copies values from inPin to all output
func Connector(name string, inPin Pin, size int) (output []Pin) {
	log.Printf("connector starting: %s", name)
	wires := make([]datapath.Wires, size)

	// connect wires to output pins
	for i := range wires {
		wires[i] = datapath.MakeWires()
		output = append(output, wires[i])
	}

	// device loop
	go func() {
		for {
			select {
			case val := <-inPin:
				log.Printf("connector:%s {%032b copied to %d pins}", name, val, size)
				for _, wire := range wires {
					wire <- val
				}
			}
		}
	}()

	return
}
