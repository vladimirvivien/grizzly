package device

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestMux(t *testing.T) {
	choices := 12
	var wires []datapath.Wires
	var muxWires []datapath.WireRcvd

	// create wires
	for i := 0; i < choices; i++ {
		w := datapath.MakeWires()
		wires = append(wires, w)
		muxWires = append(muxWires, w)
	}

	// setup mux
	selPin := datapath.MakeWires()
	muxOut := Mux(selPin, muxWires...)

	// send and test data
	for i := 0; i < choices; i++ {
		// put data on wires
		go func() {
			for i, w := range wires {
				w <- uint32(i * 2)
			}
		}()
		time.Sleep(50 * time.Millisecond) // wait for data propagation

		sel := rand.Intn(choices)
		selPin <- uint32(sel)
		t.Logf("selected mux line %d", sel)
		out := <-muxOut
		t.Logf("out: %d", out)
		if out != uint32(sel*2) {
			t.Fatalf("expecting %d, got %d", sel*2, out)
		}
	}

}

func TestFanout(t *testing.T) {
	inWire := datapath.MakeWires()
	fanCount := 4
	outWires := Fanout(inWire, fanCount)
	waiter := make(chan struct{})

	go func() {
		inWire <- 2
		inWire <- 4
		inWire <- 6
		inWire <- 2
		close(inWire)
	}()

	go func() {
		var wg sync.WaitGroup
		wg.Add(fanCount) // sync on total fanout wires
		for i, wire := range outWires {
			go func(idx int, pin Pin) {
				if val := <-pin; val != 2 {
					t.Fatalf("Unexpected value %d", val)
				}
				if val := <-pin; val != 4 {
					t.Fatalf("Unexpected value %d", val)
				}
				if val := <-pin; val != 6 {
					t.Fatalf("Unexpected value %d", val)
				}
				if val := <-pin; val != 2 {
					t.Fatalf("Unexpected value %d", val)
				}
				wg.Done()
			}(i, wire)
		}
		wg.Wait()
		close(waiter)
	}()

	select {
	case <-waiter:
	case <-time.After(5000 * time.Millisecond):
		t.Fatal("Fanout test took too long")
	}
}
