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
			var value uint32
			select {
			case sel := <-selPin:
				select {
				case value = <-inPins[sel]:
					output <- value
				}
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
