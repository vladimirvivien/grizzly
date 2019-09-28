package reg

import (
	"testing"
	"time"
)

func TestRegisterFileRS(t *testing.T) {
	regdata := []uint32{
		0: 0x0,
		1: 0x1,
		2: 0x2,
		3: 0x3,
		4: 0x4,
		5: 0x5,
		6: 0x6,
		7: 0xA,
		9: 0x9,
	}

	tests := []struct {
		name    string
		rslines func() (rs1 <-chan uint32, rs2 <-chan uint32)
		evalrs  func(t *testing.T, rf *RegisterFile)
	}{
		{
			name: "test rs values",
			rslines: func() (rs1, rs2 <-chan uint32) {
				rs1wire := make(chan uint32)
				go func() {
					rs1wire <- 0x0
					rs1wire <- 0x1
					close(rs1wire)
				}()

				rs2wire := make(chan uint32)
				go func() {
					rs2wire <- 0x2
					rs2wire <- 0x3
					rs2wire <- 0x4
					close(rs2wire)
				}()
				return rs1wire, rs2wire
			},
			evalrs: func(t *testing.T, rf *RegisterFile) {
				rf.Load()
				rs1 := <-rf.RS1DataOut()
				if rs1 != 0x0 {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs2 := <-rf.RS2DataOut()
				if rs2 != 0x2 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}

				rf.Load()
				rs1 = <-rf.RS1DataOut()
				if rs1 != 0x1 {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs2 = <-rf.RS2DataOut()
				if rs2 != 0x3 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}

				rf.Load()
				rs1 = <-rf.RS1DataOut()
				if rs1 != 0x0 { // closed
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs2 = <-rf.RS2DataOut()
				if rs2 != 0x4 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})
			reg := newRegister()
			reg.file = regdata

			// wire ports
			rs1wire, rs2wire := test.rslines()
			reg.RS1AddrIn(rs1wire)
			reg.RS2AddrIn(rs2wire)

			// start component
			if err := reg.Run(); err != nil {
				t.Fatal(err)
			}

			// setup evaluation
			go func() {
				defer close(wait)
				test.evalrs(t, reg)
			}()

			// detect stuck path
			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatal("Register operation took too long to comlete...")
			}
		})
	}
}

func TestRegisterFileRD(t *testing.T) {
	regdata := []uint32{
		0: 0xCAFFE,
		1: 0xCAFFE,
		2: 0xCAFFE,
		3: 0xCAFFE,
		4: 0xCAFFE,
		5: 0xCAFFE,
		6: 0xCAFFE,
		7: 0xCAFFE,
		9: 0xCAFFE,
	}

	tests := []struct {
		name    string
		rdlines func() (rd <-chan uint32, data <-chan uint32)
		evalrd  func(t *testing.T, rf *RegisterFile)
	}{
		{
			name: "test rd values",
			rdlines: func() (<-chan uint32, <-chan uint32) {
				rd := make(chan uint32)
				go func() {
					rd <- 0x5
					rd <- 0x9
					rd <- 0x9 // overwrites last value
					close(rd)
				}()
				data := make(chan uint32)
				go func() {
					data <- 0xAB
					data <- 0xF1 // overwritten
					data <- 0xA1
					close(data)
				}()

				return rd, data
			},
			evalrd: func(t *testing.T, rf *RegisterFile) {
				rs1 := make(chan uint32)
				rf.rs1AddrIn = rs1
				go func() {
					rs1 <- 0x5
					rs1 <- 0x9
					rs1 <- 0x9
					close(rs1)
				}()

				rf.Store() // 0x5
				rf.Load()  // 0xAB
				data := <-rf.RS1DataOut()
				if data != 0xAB {
					t.Errorf("RS1 has unexpected data: %0x", data)
				}

				rf.Store() // 0x9
				rf.Load()  // 0xF1
				data = <-rf.RS1DataOut()
				if data != 0xF1 {
					t.Errorf("RS1 has unexpected data: %0x", data)
				}

				rf.Store() // 0x9
				rf.Load()  // 0xA1
				data = <-rf.RS1DataOut()
				if data != 0xA1 {
					t.Errorf("RS1 has unexpected data: %0x", data)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := make(chan struct{})
			reg := newRegister()
			reg.file = regdata

			// wire ports
			rd, data := test.rdlines()
			reg.RDAddrIn(rd)
			reg.DataIn(data)

			// start component
			if err := reg.Run(); err != nil {
				t.Fatal(err)
			}

			// setup evaluation
			go func() {
				defer close(wait)
				test.evalrd(t, reg)
			}()

			// detect stuck path
			select {
			case <-wait:
			case <-time.After(5 * time.Millisecond):
				t.Fatal("Register operation took too long to comlete...")
			}
		})
	}
}
