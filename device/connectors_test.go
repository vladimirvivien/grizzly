package device

import (
	"sync"
	"testing"
	"time"
)

func TestFanout(t *testing.T) {
	inWire := MakeWires()
	fanCount :=4
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
				if val := <-pin; val != 2 {t.Fatalf("Unexpected value %d", val)}
				if val := <-pin; val != 4 {t.Fatalf("Unexpected value %d", val)}
				if val := <-pin; val != 6 {t.Fatalf("Unexpected value %d", val)}
				if val := <-pin; val != 2 {t.Fatalf("Unexpected value %d", val)}
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
