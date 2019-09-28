package decoder

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/inst"
)

func decode(i inst.Type) (*inst.Fields, error) {
	switch op := decodeOp(i); op {
	case inst.Opcodes.R:
		return &inst.Fields{
			Opcode: op,
			Rd:     decodeRd(1),
			Funct3: decodeFunct3(i),
			Rs1:    decodeRs1(i),
			Rs2:    decodeRs2(i),
			Funct7: decodeFunct7(i),
		}, nil
	default:
		return nil, fmt.Errorf("inst unsupported: %b", i)
	}
}

func decodeOp(i inst.Type) uint32 {
	return uint32(i) & 0x7F
}

func decodeRd(i inst.Type) uint32 {
	return (uint32(i) >> 7) & 0x1F
}

func decodeFunct3(i inst.Type) uint32 {
	return (uint32(i) >> 12) & 0x7
}

func decodeRs1(i inst.Type) uint32 {
	return (uint32(i) >> 15) & 0x1F
}

func decodeRs2(i inst.Type) uint32 {
	return (uint32(i) >> 20) & 0x1F
}

func decodeFunct7(i inst.Type) uint32 {
	return (uint32(i) >> 25) & 0x7F
}
