// +build rv64i

package isa

func GetOpcode(i uint64) Opcode {
	return  Opcode(i & 0x7F)
}
