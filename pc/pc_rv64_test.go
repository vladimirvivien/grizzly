//go:build rv64 || rv64i

package pc

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestPCNew_RV64(t *testing.T) {
	pc := New()
	if pc.counter != 0 {
		t.Errorf("unexpected PC counter value: %d", pc.counter)
	}
	if pc.GetPin(Labels.OutCounter) == nil {
		t.Errorf("output pin not set: %s", Labels.OutCounter)
	}
}

func TestPCRun_Counter_RV64(t *testing.T) {
	pc := New()
	opCh := make(chan []byte)
	pc.Connect(Labels.InPcOp, opCh)
	if err := pc.Run(); err != nil {
		t.Fatal(err)
	}

	waiter := make(chan struct{})
	maxCount := 1000
	go func() {
		for i := 0; i < maxCount; i++ {
			opCh <- datapath.EncodePcOp(datapath.PcOp{})
		}
		close(opCh)
	}()
	go func() {
		defer close(waiter)
		for stream := range pc.GetPin(Labels.OutCounter) {
			count := datapath.DecodePC(stream)
			if count >= uint64(maxCount) {
				pc.clock.Stop()
				return
			}
		}
	}()

	select {
	case <-waiter:
	case <-time.After(1 * time.Second):
		t.Fatal("Operations took too long to complete")
	}
}

func TestPCRun_Jump_RV64(t *testing.T) {
	pc := New()
	opCh := make(chan []byte)
	pc.Connect(Labels.InPcOp, opCh)
	if err := pc.Run(); err != nil {
		t.Fatal(err)
	}

	maxCount := 1000
	go func() {
		for i := 8; i <= maxCount; i += 8 {
			opCh <- datapath.EncodePcOp(datapath.PcOp{Jump: 1, PC: datapath.XWord(i)})
		}
		close(opCh)
	}()

	waiter := make(chan struct{})
	go func() {
		defer close(waiter)
		for stream := range pc.GetPin(Labels.OutCounter) {
			count := datapath.DecodePC(stream)
			if count%8 != 0 {
				t.Errorf("unexpected counter %d", count)
			}
			if count >= uint64(maxCount) {
				pc.clock.Stop()
				return
			}
		}
	}()

	select {
	case <-waiter:
	case <-time.After(1 * time.Second):
		t.Fatal("Operations took too long to complete")
	}
}

func FuzzPC(f *testing.F) {
	f.Add(uint8(0), uint64(0x123456789ABCDE00))
	f.Add(uint8(1), uint64(0x9ABCDEF012345600))
	f.Fuzz(func(t *testing.T, isJump uint8, jumpPC uint64) {
		// Enforce alignment to 4 bytes for PC
		jumpPC = jumpPC &^ 3

		pc := New()
		opCh := make(chan []byte)
		pc.Connect(Labels.InPcOp, opCh)
		if err := pc.Run(); err != nil {
			t.Fatal(err)
		}

		go func() {
			opCh <- datapath.EncodePcOp(datapath.PcOp{
				Jump: isJump % 2,
				PC:   datapath.XWord(jumpPC),
			})
			close(opCh)
		}()

		var expected uint64
		// PC starts at 0. First output is 0.
		// Then it processes the pcOp.
		if (isJump % 2) > 0 {
			expected = jumpPC
		} else {
			expected = 4
		}

		outCh := pc.GetPin(Labels.OutCounter)
		
		// First pc.out value is the initial 0
		<-outCh

		// Second value is the updated PC value after processing pcOp
		stream, ok := <-outCh
		if !ok {
			t.Fatal("expected second output from PC")
		}
		
		pc.clock.Stop()

		actual := datapath.DecodePC(stream)
		if actual != expected {
			t.Errorf("got %x, expected %x", actual, expected)
		}
	})
}
