//go:build rv32 || rv32i || (!rv64 && !rv64i && !rv128)

package core

import (
	"encoding/binary"
	"io/ioutil"
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
)

func TestCore_Run_Manual(t *testing.T) {
	var program []byte
	program = append(program, datapath.EncodeXWord(0b000000000010_00000_000_00001_0010011)...)  // addi x1, x0, 2;  load x1=2
	program = append(program, datapath.EncodeXWord(0b000000000100_00000_000_00010_0010011)...)  // addi x2, x0, 4;  load x2=4
	program = append(program, datapath.EncodeXWord(0b000000001100_00000_000_00111_0010011)...)  // addi x7, x0, 12; load x7=12
	program = append(program, datapath.EncodeXWord(0b000000000100_00001_000_00101_0010011)...)  // addi x5, x1, 4;  x5=6
	program = append(program, datapath.EncodeXWord(0b0000000_00010_00111_000_00011_0110011)...) // add  x3, x7, x2; x3=16
	program = append(program, datapath.EncodeXWord(0b000000000001_00010_001_00110_0010011)...)  // slli x6, x2, 1;  x6=8

	t.Log("prog size:", len(program))
	cor := New()
	cor.imem.SetStore(program)

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	// stall to wait for all instructions before assessment
	<-time.After(20 * time.Millisecond)
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

// Test
//
//addi x20, x0, 2
//addi x21, x0, 4
//addi x22, x0, 12
//addi x24, x20, 4
//add  x25, x21, x22
//slli x26, x21, 1
//
//func TestCore_Run_Stream(t *testing.T) {
//	stream, err := coretest.StreamFromFile("../testing/programs/rtypes_rv32/add.bin")
//	if err != nil {
//		t.Fatal(err)
//	}
//	cor := New()
//	cor.Input(stream)
//	if err := cor.Run(); err != nil {
//		t.Fatal(err)
//	}
//
//	// stall to wait for all instructions before assessment
//	<-clock.New(300 * time.Microsecond).Ticks()
//	val := cor.reg.Probe(20)
//	if val != 2 {
//		t.Errorf("unexpected register value: x20=%d", val)
//	}
//	val = cor.reg.Probe(21)
//	if val != 4 {
//		t.Errorf("unexpected register value: x21=%d", val)
//	}
//	val = cor.reg.Probe(22)
//	if val != 12 {
//		t.Errorf("unexpected register value: x22=%d", val)
//	}
//
//	val = cor.reg.Probe(24)
//	if val != 6 {
//		t.Errorf("unexpected register value: x24=%d", val)
//	}
//	val = cor.reg.Probe(25)
//	if val != 16 {
//		t.Errorf("unexpected register value: x25=%d", val)
//	}
//	val = cor.reg.Probe(26)
//	if val != 8 {
//		t.Errorf("unexpected register value: x26=%d", val)
//	}
//}
//
//func BenchmarkCore_Run(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		stream, err := coretest.StreamFromFile("../testing/programs/rtypes_rv32/add.bin")
//		if err != nil {
//			b.Fatal(err)
//		}
//		cor := New()
//		cor.Input(stream)
//		if err := cor.Run(); err != nil {
//			b.Fatal(err)
//		}
//
//		// stall to wait for all instructions before assessment
//		<-time.After(time.Millisecond)
//		val := cor.reg.Probe(20)
//		if val != 2 {
//			b.Errorf("unexpected register value: x20=%d", val)
//		}
//		val = cor.reg.Probe(21)
//		if val != 4 {
//			b.Errorf("unexpected register value: x21=%d", val)
//		}
//		val = cor.reg.Probe(22)
//		if val != 12 {
//			b.Errorf("unexpected register value: x22=%d", val)
//		}
//
//		val = cor.reg.Probe(24)
//		if val != 6 {
//			b.Errorf("unexpected register value: x24=%d", val)
//		}
//		val = cor.reg.Probe(25)
//		if val != 16 {
//			b.Errorf("unexpected register value: x25=%d", val)
//		}
//		val = cor.reg.Probe(26)
//		if val != 8 {
//			b.Errorf("unexpected register value: x26=%d", val)
//		}
//	}
//}

func instToStream(word datapath.XWord) []byte {
	inst := make([]byte, 4)
	binary.LittleEndian.PutUint32(inst, word)
	return inst
}

func TestCore_Run_Binary(t *testing.T) {
	// Load the binary program file compiled by Zig
	content, err := ioutil.ReadFile("../testing/programs/rtypes_rv32/build/add.bin")
	if err != nil {
		t.Fatal(err)
	}

	cor := New()
	cor.imem.SetStore(content)

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	// wait for execution
	<-time.After(20 * time.Millisecond)

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

func TestCore_Run_Branch(t *testing.T) {
	// Load the branch test binary file compiled by Zig
	content, err := ioutil.ReadFile("../testing/programs/branch_rv32/build/branch_test.bin")
	if err != nil {
		t.Fatal(err)
	}

	cor := New()
	cor.imem.SetStore(content)

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	// wait for execution
	<-time.After(20 * time.Millisecond)

	// Check register values
	val := cor.reg.Probe(1) // x1 = 5
	if val != 5 {
		t.Errorf("unexpected register value: x1=%d", val)
	}
	val = cor.reg.Probe(2) // x2 = 5
	if val != 5 {
		t.Errorf("unexpected register value: x2=%d", val)
	}
	val = cor.reg.Probe(3) // x3 = 0 (should not be executed/modified to 1 or 2)
	if val != 0 {
		t.Errorf("unexpected register value: x3=%d (branch was not taken correctly)", val)
	}
	val = cor.reg.Probe(4) // x4 = 10
	if val != 10 {
		t.Errorf("unexpected register value: x4=%d", val)
	}
	val = cor.reg.Probe(5) // x5 = 6
	if val != 6 {
		t.Errorf("unexpected register value: x5=%d", val)
	}
	val = cor.reg.Probe(6) // x6 = 15
	if val != 15 {
		t.Errorf("unexpected register value: x6=%d", val)
	}
}

func TestCore_Run_LoadStore(t *testing.T) {
	var program []byte
	appendInst := func(inst uint32) {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, inst)
		program = append(program, buf...)
	}

	// addi x1, x0, 256;  load address x1=256
	appendInst(0b000100000000_00000_000_00001_0010011)
	// addi x2, x0, 127;  load value x2=127
	appendInst(0b000001111111_00000_000_00010_0010011)
	// sw x2, 0(x1);      store value to memory
	appendInst(0b0000000_00010_00001_010_00000_0100011)
	// addi x3, x0, 0;    clear register x3
	appendInst(0b000000000000_00000_000_00011_0010011)
	// lw x3, 0(x1);      load value from memory
	appendInst(0b000000000000_00001_010_00011_0000011)

	cor := New()
	cor.imem.SetStore(program)

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	// wait for execution
	<-time.After(20 * time.Millisecond)

	val := cor.reg.Probe(3)
	if val != 127 {
		t.Errorf("expected loaded value to be 127, got %d", val)
	}
}

func TestCore_Run_Branch_Manual(t *testing.T) {
	var program []byte
	appendInst := func(inst uint32) {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, inst)
		program = append(program, buf...)
	}

	// addi x1, x0, 10
	appendInst(0b000000001010_00000_000_00001_0010011)
	// addi x2, x0, 10
	appendInst(0b000000001010_00000_000_00010_0010011)
	// beq x1, x2, 12 (offset of +12 bytes skips the next 2 addi instructions)
	appendInst(0b0000000_00010_00001_000_0110_0_1100011)
	// addi x3, x0, 1 (skipped)
	appendInst(0b000000000001_00000_000_00011_0010011)
	// addi x3, x0, 2 (skipped)
	appendInst(0b000000000010_00000_000_00011_0010011)
	// addi x4, x0, 100 (target instruction)
	appendInst(0b000001100100_00000_000_00100_0010011)

	cor := New()
	cor.imem.SetStore(program)

	if err := cor.Run(); err != nil {
		t.Fatal(err)
	}

	// wait for execution
	<-time.After(20 * time.Millisecond)

	x3 := cor.reg.Probe(3)
	if x3 != 0 {
		t.Errorf("expected x3 to remain 0, got %d", x3)
	}

	x4 := cor.reg.Probe(4)
	if x4 != 100 {
		t.Errorf("expected target x4 to be 100, got %d", x4)
	}
}

