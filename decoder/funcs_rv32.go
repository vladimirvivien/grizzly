package decoder

import (
	"encoding/binary"

	"github.com/vladimirvivien/grizzly/datapath"
)

// bytesToInst converts a group of 4 bytes to
// instruction Word
func bytesToInst(bits []byte) datapath.XWord {
	return binary.LittleEndian.Uint32(bits)
}
