//go:build rv64 || rv64i

package load

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

// Decode 64-bit load instruction format
func Decode(i datapath.XWord) datapath.OpFields {
	var fields datapath.OpFields
	fields.Opcode = uint8(i & 0x7F)
	fields.Rd = uint8((i >> 7) & 0x1F)
	fields.Funct3 = uint8((i >> 12) & 0x7)
	fields.Rs1 = uint8((i >> 15) & 0x1F)
	val := (i >> 20) & 0xFFF
	if (val & 0x800) != 0 {
		val |= 0xfffff000
	}
	fields.Imm = uint32(val)
	return fields
}
