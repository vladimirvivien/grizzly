//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package pc

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestPCNew(t *testing.T) {
	pc := New()
	if pc.counter != 0 {
		t.Errorf("unpexpected PC counter value: %d", pc.counter)
	}
	if pc.GetPin(Labels.OutCounter) == nil {
		t.Errorf("output pin not set: %s", Labels.OutCounter)
	}
}

func TestPCRun_Counter(t *testing.T) {
	pc := New()
	opCh := make(chan []byte)
	pc.Connect(Labels.InPcOp, opCh)
	if err := pc.Run(); err != nil {
		t.Fatal(err)
	}

	waiter := make(chan struct{})
	maxCount := 1024 * 100
	go func() {
		for i := 0; i < maxCount; i++ {
			// jump = 0; means transfer pc = pc + 4
			opCh <- datapath.EncodePcOp(datapath.PcOp{})
		}
		close(opCh)
	}()
	go func() {
		defer close(waiter)
		for stream := range pc.GetPin(Labels.OutCounter) {
			count := datapath.DecodePC(stream)
			if count >= uint32(maxCount) {
				pc.clock.Stop()
				return
			}
		}
	}()

	select {
	case <-waiter:
	case <-time.After(time.Second):
		t.Fatal("Operations took too long to complete")
	}
}

func TestPCRun_Jump(t *testing.T) {
	pc := New()
	opCh := make(chan []byte)
	pc.Connect(Labels.InPcOp, opCh)
	if err := pc.Run(); err != nil {
		t.Fatal(err)
	}

	maxCount := 1024 * 100
	go func() {
		for i := 8; i <= maxCount; i+=8 {
			opCh <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC:datapath.XWord(i)})
		}
		close(opCh)
	}()

	waiter := make(chan struct{})
	go func() {
		defer close(waiter)
		for stream := range pc.GetPin(Labels.OutCounter) {
			count := datapath.DecodePC(stream)
			if count % 8 != 0 {
				t.Errorf("unexpected counter %d", count)
			}
			if count >= uint32(maxCount) {
				pc.clock.Stop()
				return
			}
		}
	}()

	select {
	case <-waiter:
	case <-time.After(time.Second):
		t.Fatal("Operations took too long to complete")
	}
}