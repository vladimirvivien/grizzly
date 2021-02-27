// +build rv32i

package isa

func GetOpcode(i uint32) Opcode {
	return  Opcode(i & 0x7F)
}
