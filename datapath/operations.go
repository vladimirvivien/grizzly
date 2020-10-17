package datapath

type Packet struct {
	Word
	Wires
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

// Collect collects data (word) serially from wires
// and retries until all wires successfully return data
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
