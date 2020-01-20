package decoder

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/isa"
)

func decode(i isa.Inst) (*isa.Fields, error) {
	switch op := decodeOp(i); op {
	case isa.Opcodes.R:
		return &isa.Fields{
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

func decodeOp(i isa.Inst) uint32 {
	return uint32(i) & 0x7F
}

func decodeRd(i isa.Inst) uint32 {
	return (uint32(i) >> 7) & 0x1F
}

func decodeFunct3(i isa.Inst) uint32 {
	return (uint32(i) >> 12) & 0x7
}

func decodeRs1(i isa.Inst) uint32 {
	return (uint32(i) >> 15) & 0x1F
}

func decodeRs2(i isa.Inst) uint32 {
	return (uint32(i) >> 20) & 0x1F
}

func decodeFunct7(i isa.Inst) uint32 {
	return (uint32(i) >> 25) & 0x7F
}
