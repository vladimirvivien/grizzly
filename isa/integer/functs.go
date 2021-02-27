package integer

// Functs encodes integer R-format functs fields funct7 and funct3
// as a single value into the lower 10 bits of the result:
//
//    [00000077 77777333]
//
func Functs(f7, f3 uint8) (result uint16) {
	result = (result | uint16(f7)) << 3
	result = result | uint16(f3)
	return result
}
