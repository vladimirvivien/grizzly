//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package testing

import (
	"testing"
)

func TestLoadBuffer(t *testing.T) {
	prog := "./programs/rtypes_rv32/add.bin"
	mem, err := LoadFile(prog)
	if err != nil {
		t.Fatal(err)
	}
	PrintBinary(mem)
}