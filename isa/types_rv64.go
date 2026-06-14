//go:build rv64 || rv64i

package isa

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

func GetOpcode(i datapath.XWord) Opcode {
	return  Opcode(i & 0x7F)
}
