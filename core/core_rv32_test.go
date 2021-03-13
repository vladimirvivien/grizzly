package core

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/clock"
	"github.com/vladimirvivien/grizzly/datapath"
	coretest "github.com/vladimirvivien/grizzly/testing"
)

func TestCore_Run_Manual(t *testing.T) {
	cor := New()
	ch := make(chan []byte)
	cor.Input(ch)
	go func(){
		// data load
		ch <- instToStream(0b000000000010_00000_000_00001_0010011)  // addi x1, x0, 2;  load x1=2
		ch <- instToStream(0b000000000100_00000_000_00010_0010011)  // addi x2, x0, 4;  load x2=4
		ch <- instToStream(0b000000001100_00000_000_00111_0010011)  // addi x7, x0, 12; load x7=12

		ch <- instToStream(0b000000000100_00001_000_00101_0010011)  // addi x5, x1, 4;  x5=6
		ch <- instToStream(0b0000000_00010_00111_000_00011_0110011) // add  x3, x7, x2; x3=16
		ch <- instToStream(0b000000000001_00010_001_00110_0010011)  // slli x6, x2, 1;  x6=8
		close(ch)
	}()

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	// stall to wait for all instructions before assessment
	<-clock.New(100*time.Microsecond).Ticks()
	val := cor.reg.Probe(1)
	if val!= 2{
		t.Errorf("unexpected register value: x1=%d", val)
	}
	val = cor.reg.Probe(2)
	if val!= 4{
		t.Errorf("unexpected register value: x2=%d", val)
	}
	val = cor.reg.Probe(7)
	if val!= 12{
		t.Errorf("unexpected register value: x7=%d", val)
	}
	val = cor.reg.Probe(5)
	if val!= 6{
		t.Errorf("unexpected register value: x5=%d", val)
	}
	val = cor.reg.Probe(3)
	if val!= 16{
		t.Errorf("unexpected register value: x3=%d", val)
	}
	val = cor.reg.Probe(6)
	if val!= 8{
		t.Errorf("unexpected register value: reg[6]=%d", val)
	}
}


func TestCore_Run_Stream(t *testing.T) {
	stream, err := coretest.StreamFromFile("../testing/programs/rtypes_rv32/add.bin")
	if err != nil{
		t.Fatal(err)
	}
	cor := New()
	cor.Input(stream)
	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	// stall to wait for all instructions before assessment
	//<-clock.New(100*time.Microsecond).Ticks()
	//val := cor.reg.Probe(20)
	//if val!= 2{
	//	t.Errorf("unexpected register value: x1=%d", val)
	//}
}

func instToStream(word datapath.XWord)[]byte {
	inst := make([]byte,4)
	binary.LittleEndian.PutUint32(inst, word)
	return inst
}