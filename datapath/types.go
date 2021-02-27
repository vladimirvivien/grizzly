package datapath

// TODO
// Investigate using []byte to represent instruction streams between components. This
// approach would allow Grizzly to support multi-width instructions (32, 64, 128, etc).
//
// How would this work?
//   - The Wires type would be defined as chan []byte to stream both control/program data.
//   - Instructions would be converted from a numeric value to a stream of bytes
//   - The stream of bytes would be sent to core components over Wires
//   - When components get []byte, the bytes are converted to numeric representation
//     using the encoding/binary package to narrow to a specific value based on configured
//     instruction width.
//
// See https://play.golang.org/p/py_Uv9zSXWv
//
// This change would allow Grizzly to support different implementations of the RISCV
// ISA including compressed instructions based on the size of XLEN.
const XLEN = 32
const Width32 = 32
const Width64 = 64

// type Word = [Xlen]byte

type Word = uint32
type DoubleWord = uint64

type Wires = chan Word
type WireRcvd = <-chan Word
type WireSend = chan <- Word
