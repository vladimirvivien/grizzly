//go:build rv64 || rv64i

package core

import (
	"encoding/binary"
	"io/ioutil"
	"testing"
	"time"
)

func TestCore_Run_Manual(t *testing.T) {
	var program []byte
	appendInst := func(inst uint32) {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, inst)
		program = append(program, buf...)
	}

	appendInst(0b000000000010_00000_000_00001_0010011)  // addi x1, x0, 2;  load x1=2
	appendInst(0b000000000100_00000_000_00010_0010011)  // addi x2, x0, 4;  load x2=4
	appendInst(0b000000001100_00000_000_00111_0010011)  // addi x7, x0, 12; load x7=12
	appendInst(0b000000000100_00001_000_00101_0010011)  // addi x5, x1, 4;  x5=6
	appendInst(0b0000000_00010_00111_000_00011_0110011) // add  x3, x7, x2; x3=16
	appendInst(0b000000000001_00010_001_00110_0010011)  // slli x6, x2, 1;  x6=8

	t.Log("prog size:", len(program))
	cor := New()
	cor.imem.SetStore(program)

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Millisecond)
	val := cor.reg.Probe(1)
	if val != 2 {
		t.Errorf("unexpected register value: x1=%d", val)
	}
	val = cor.reg.Probe(2)
	if val != 4 {
		t.Errorf("unexpected register value: x2=%d", val)
	}
	val = cor.reg.Probe(7)
	if val != 12 {
		t.Errorf("unexpected register value: x7=%d", val)
	}
	val = cor.reg.Probe(5)
	if val != 6 {
		t.Errorf("unexpected register value: x5=%d", val)
	}
	val = cor.reg.Probe(3)
	if val != 16 {
		t.Errorf("unexpected register value: x3=%d", val)
	}
	val = cor.reg.Probe(6)
	if val != 8 {
		t.Errorf("unexpected register value: reg[6]=%d", val)
	}
}

func TestCore_Run_Binary(t *testing.T) {
	content, err := ioutil.ReadFile("../testing/programs/rtypes_rv64/build/add.bin")
	if err != nil {
		t.Fatal(err)
	}

	cor := New()
	cor.imem.SetStore(content)

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Millisecond)

	val := cor.reg.Probe(20)
	if val != 2 {
		t.Errorf("unexpected register value: x20=%d", val)
	}
	val = cor.reg.Probe(21)
	if val != 4 {
		t.Errorf("unexpected register value: x21=%d", val)
	}
	val = cor.reg.Probe(22)
	if val != 12 {
		t.Errorf("unexpected register value: x22=%d", val)
	}
	val = cor.reg.Probe(24)
	if val != 6 {
		t.Errorf("unexpected register value: x24=%d", val)
	}
	val = cor.reg.Probe(25)
	if val != 16 {
		t.Errorf("unexpected register value: x25=%d", val)
	}
	val = cor.reg.Probe(26)
	if val != 8 {
		t.Errorf("unexpected register value: x26=%d", val)
	}
}
