package datapath

//type Packet struct {
//	XWord
//	Wires
//}

//// MakeWires creates virtual wires to carry XWord-size data/control.
//// Wires are blocking constructs, meaning once data is sent, it must
//// be consumed by receiver before next send operation.
//// Note: wires are 1-sized buffered Go channels.
//func MakeWires() Wires {
//	return make(chan XWord,1)
//}
//
//// Send sends data (word) to wires in serial order
//// and blocks on each send until received.
//// TODO possible deadline to avoid long waits
//func Send(packets ...Packet) {
//	for _, p := range packets {
//		select {
//		case p.Wires <- p.XWord:
//		}
//	}
//}
//
//type Receiver struct {
//	Name string
//}
//
//func NewReceiver(name string) *Receiver{
//	return &Receiver{Name:name}
//}
//
//// R triggers the receive operation for the receiver component
//func (c *Receiver) R(wires ...WireRcvd)(words []XWord) {
//	return Collect(c.Name, wires...)
//}
//
//// Collect collects all data (word) serially from wires
//// and waits for each wires to be ready
//// TODO possible deadline to avoid long waits
//func Collect(collector string, wires ...WireRcvd) (words []XWord) {
//	for i, wire := range wires {
//		log.Printf("%s:receiving wire[%d]", collector, i)
//		select {
//		case word := <-wire:
//			log.Printf("%s: rcvd: wire[%d]=%032b", collector, i, word)
//			words = append(words, word)
//		}
//	}
//	return
//}
