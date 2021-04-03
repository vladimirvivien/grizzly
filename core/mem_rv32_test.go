package core

// main:
// addi x20, x0, 2
// addi x21, x0, 4
// addi x22, x0, 16
// add  x25, x20, x21
// sw   x25, 0(x22)
//func TestCore_Mem(t *testing.T) {
//	stream, err := coretest.StreamFromFile("../testing/programs/mem_rv32/dmem.bin")
//	if err != nil{
//		t.Fatal(err)
//	}
//	cor := New()
//	cor.Input(stream)
//	if err := cor.Run(); err != nil {
//		t.Fatal(err)
//	}
//
//	// stall to wait for all instructions before assessment
//	<-time.After(50*time.Millisecond)
//	val := cor.dmem.TestProbe(16)
//	if val!= 6{
//		t.Errorf("unexpected register value: M[16]=%d", val)
//	}
//}
