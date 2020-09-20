package device

func Mux(selPin Pin, inPins...Pin) Pin {
	output := MakeWires()
	go func(){
		defer close(output)
		for {
			var value uint32
			select{
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