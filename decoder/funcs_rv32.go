package decoder

import (
	"encoding/binary"

	"github.com/vladimirvivien/grizzly/datapath"
)

func bytesToInst(bits []byte) datapath.XWord {
	return binary.LittleEndian.Uint32(bits)
}
