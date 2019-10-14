package reg

import (
	"testing"
	"time"
)

func TestRegisterFile(t *testing.T) {
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
		rslines func() (rs1, rs2, data, rd <-chan uint32)
		evalrs  func(t *testing.T, rf *RegisterFile)
	}{
		{
			name: "test rs values",
			rslines: func() (rs1, rs2, data, rd <-chan uint32) {
				rs1wire := make(chan uint32)
				rs2wire := make(chan uint32)
				datawire := make(chan uint32)
				rdwire := make(chan uint32)
				go func() {
					defer func() {
						close(rs1wire)
						close(rs2wire)
					}()
					rs1wire <- 0x0
					rs1wire <- 0x1
					rs2wire <- 0x2
					rs2wire <- 0x3
					rs2wire <- 0x4

					//write data
					// data must always be specified before rd
					// or risk deadlock
					datawire <- 0xCAFE
					rdwire <- 0x05
					rs2wire <- 0x5 // read the data

					rs1wire <- 0x7
					rs2wire <- 0x5 // reread
				}()
				return rs1wire, rs2wire, datawire, rdwire
			},
			evalrs: func(t *testing.T, rf *RegisterFile) {
				//rf.Enable()
				rs1 := <-rf.RS1DataOut()
				if rs1 != 0x0 {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs1 = <-rf.RS1DataOut()
				if rs1 != 0x1 {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs2 := <-rf.RS2DataOut()
				if rs2 != 0x2 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs2 = <-rf.RS2DataOut()
				if rs2 != 0x3 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs2 = <-rf.RS2DataOut()
				if rs2 != 0x4 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs2 = <-rf.RS2DataOut()
				if rs2 != 0xCAFE {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs1 = <-rf.RS1DataOut()
				if rs1 != 0xA {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs2 = <-rf.RS2DataOut()
				if rs2 != 0xCAFE {
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
			rs1wire, rs2wire, data, rd := test.rslines()
			reg.RS1AddrIn(rs1wire)
			reg.RS2AddrIn(rs2wire)
			reg.DataIn(data)
			reg.RDAddrIn(rd)

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
