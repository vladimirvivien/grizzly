package device

/*
Selectors work similarly as do a multiplexer, however, without a select line.
They are constructed using two or more Pins (Go channels) as input.
Initially, all receive operations are blocked until a Pin receives a
value.  The selector will automatically select and return the received
value on its output Pin and then blocked until the next value.
 */

// Select2 builds a selector with two input pins
func Select2(in0, in1 Pin) Pin {
	output := MakeWires()
	go func() {
		defer close(output)
		for {
			var value uint32
			select {
			case val := <-in0:
				value = val
			case val := <-in1:
				value = val
			}
			output <- value
		}
	}()

	return output
}

// Select4 builds a selector with 4 inputs
func Select4(in0, in1, in2, in3 Pin) Pin {
	output := MakeWires()
	go func() {
		defer close(output)
		for {
			var value uint32
			select {
			case val := <-in0:
				value = val
			case val := <-in1:
				value = val
			case val := <-in2:
				value = val
			case val := <-in3:
				value = val
			}
			output <- value
		}
	}()

	return output
}