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