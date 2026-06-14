//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package isa

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

func GetOpcode(i datapath.XWord) Opcode {
	return  Opcode(i & 0x7F)
}
