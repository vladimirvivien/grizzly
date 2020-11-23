package datapath

type Packet struct {
	Word
	Wires
}
func MakeWires() Wires {
	return make(chan Word)
}

// Send sends data (word) to wires in serial order
// and blocks on each send until received.
// TODO possible deadline to avoid long waits
func Send(packets ...Packet) {
	for _, p := range packets {
		select {
		case p.Wires <- p.Word:
		}
	}
}

// Collect collects all data (word) serially from wires
// and waits for each wires to be ready
// TODO possible deadline to avoid long waits
func Collect(wires ...WireRcvd) (words []Word) {
	for _, wire := range wires {
		select {
		case word := <-wire:
			words = append(words, word)
		}
	}
	return
}