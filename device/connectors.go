package device

// Mux uses selector pin as index to select value from
// a collection of pins (channels)
func Mux(selPin Pin, inPins ...Pin) Pin {
	output := MakeWires()
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
	wires := make([]Wires, fanSize)

	// connect wires to output pins
	for i := range wires {
		wires[i] = MakeWires()
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

// TODO this may be needed
// func FanoutSequential(inPin Pin, wires...Wires)
