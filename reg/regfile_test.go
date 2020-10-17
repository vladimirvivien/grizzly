package reg

import (
	"testing"
	"time"

	"github.com/vladimirvivien/grizzly/datapath"
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
		rslines func() (rs1, rs2, data, rd, werf <-chan uint32)
		evalrs  func(t *testing.T, rf *RegisterFile)
	}{
		{
			name: "test rs values",
			rslines: func() (rs1, rs2, data, rd, werf <-chan uint32) {
				rs1wire := datapath.MakeWires()
				rs2wire := datapath.MakeWires()
				datawire := datapath.MakeWires()
				rdwire := datapath.MakeWires()
				wenable := datapath.MakeWires()
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
					// be mindful of order of rd and rdData.
					wenable <- 1
					rdwire <- 0x05
					datawire <- 0xCAFE
					rs2wire <- 0x5 // read the data

					rs1wire <- 0x7
					rs2wire <- 0x5 // reread
				}()
				return rs1wire, rs2wire, datawire, rdwire, wenable
			},
			evalrs: func(t *testing.T, rf *RegisterFile) {
				rs1 := <-rf.GetPin(Out.RS1Data)
				if rs1 != 0x0 {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs1 = <-rf.GetPin(Out.RS1Data)
				if rs1 != 0x1 {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs2 := <-rf.GetPin(Out.RS2Data)
				if rs2 != 0x2 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs2 = <-rf.GetPin(Out.RS2Data)
				if rs2 != 0x3 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs2 = <-rf.GetPin(Out.RS2Data)
				if rs2 != 0x4 {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs2 = <-rf.GetPin(Out.RS2Data)
				if rs2 != 0xCAFE {
					t.Errorf("Unexpected RS2 data %d", rs2)
				}
				rs1 = <-rf.GetPin(Out.RS1Data)
				if rs1 != 0xA {
					t.Errorf("Unexpected RS1 data %d", rs1)
				}
				rs2 = <-rf.GetPin(Out.RS2Data)
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
			rs1wire, rs2wire, data, rd, werf := test.rslines()
			reg.SetPin(In.RS1Addr, rs1wire)
			reg.SetPin(In.RS2Addr, rs2wire)
			reg.SetPin(In.Data, data)
			reg.SetPin(In.RDAddr, rd)
			reg.SetPin(In.Werf, werf)

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
