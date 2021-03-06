package integer

// Functs encodes integer R-format functs fields funct7 and funct3
// as a single value into the lower 10 bits of the result:
//
//    [XXXXXX77 77777333]
//
func Functs(f7, f3 uint8) (functs uint16) {
	functs = (functs | uint16(f7)) << 3
	functs = functs | uint16(f3)
	return
}

// Defuncts extracts ISA function fields funct7 and funct3 assuming
// value functs contain these values concatenated in the lower 10 bits as:
//
//    [XXXXXX77 77777333]
//
func Defuncts(functs uint16) (funct7, funct3 uint8) {
	funct3 = uint8(functs & 0x7)
	funct7 = uint8((functs >> 3) & 0x7F)
	return
}
