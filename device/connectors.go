package device

import (
	"log"

	"github.com/vladimirvivien/grizzly/datapath"
)

// Mux uses selector pin as index to select value from
// a collection of pins (channels)
func Mux(selPin Pin, inPins ...Pin) Pin {
	log.Printf("mux: 1-to-%d created", len(inPins))
	output := datapath.MakeWires()
	go func() {
		defer close(output)
		for {
			words := datapath.Collect(inPins...)
			log.Printf("mux: collected inputs: %d", len(words))
			select {
			case sel := <-selPin:
				val := words[sel]
				log.Printf("mux: {sel:%d, lines:%d, value:%032b}", sel, len(inPins), val)
				output <- val
			}
		}
	}()
	return output
}

// Fanout transfer value read from inPin to output
func Fanout(inPin Pin, fanSize int) (output []Pin) {
	wires := make([]datapath.Wires, fanSize)

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
				for _, wire := range wires {
					wire <- val
				}
			}
		}
	}()

	return
}
