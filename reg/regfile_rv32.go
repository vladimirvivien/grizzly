package reg

import (
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

var(
	Labels = struct {
		InFields datapath.Pin
		InData datapath.Pin
		OutAluParams datapath.Pin
	}{
		InFields: datapath.Pin("regfile.in.opfields"),
		InData: datapath.Pin("regfile.in.data"),
		OutAluParams: datapath.Pin("regfile.out.aluparams"),
	}
)

type writeSignal = struct{}
type regfile = []datapath.XWord
type RegisterFile struct {
	*datapath.BaseComponent
	m         sync.RWMutex
	file      regfile
	writeSig  chan writeSignal
	output    chan []byte
}

func New() *RegisterFile {
	reg := &RegisterFile{
		BaseComponent: datapath.NewBase(),
		writeSig: make(chan writeSignal),
		file:     make(regfile, datapath.RegSize, datapath.RegSize),
		output:   make(chan []byte),
	}
	reg.Connect(Labels.OutAluParams, reg.output)
	return reg
}


// Run starts the register file component
func (r RegisterFile) Run() error {
	input := r.GetPin(Labels.InFields)
	if input == nil {
		return fmt.Errorf("register file: missing input: %s", Labels.InFields)
	}
	inData := r.GetPin(Labels.InData)
	if inData == nil {
		return fmt.Errorf("register file: missing data input: %s", Labels.InData)
	}

	// Register instruction input loop
	// This loop handles operation fields from the decoder.
	// Once processed, the register will setup control operations
	// and prepare operands for the ALU.
	// A semaphore signal is used to wait for writebacks
	// from operations (Op.R, Op.RI, Op.L, etc) which requires it.
	go func() {
		defer close(r.output)
		for stream := range input {
			op := datapath.DecodeOpFields(stream)
			params := datapath.Operation{
				Opcode: op.Opcode,
				Rd:     op.Rd,
				Funct3: op.Funct3,
				Funct7: op.Funct7,
			}

			// Select ALU operands:
			// Select between register data
			// or immediate values to send to
			// ALU.
			switch op.Opcode {
			case isa.Opcodes.R:
				params.Op1 = r.read(op.Rs1)
				params.Op2 = r.read(op.Rs2)

				// write output,
				// wait for writeback signal before next read
				r.output <- datapath.EncodeOp(params)
				<-r.writeSig

			case isa.Opcodes.RI:
				params.Op1 = r.read(op.Rs1)
				switch op.Funct3 {
				case integer.Slli.F3, integer.Srli.F3, integer.Srai.F3:
					params.Op2 = datapath.XWord(op.Shift)
				default:
					params.Op2 = op.Imm
				}

				r.output <- datapath.EncodeOp(params)
				<-r.writeSig
			}
		}
	}()

	// Register data input loop
	// Handles register data storage request from downsream operations.
	// A semaphore signal is used to ensure that the read-loop can only
	// proceed after a previous write.
	go func() {
		for dataStream := range inData {
			data := datapath.DecodeRegStore(dataStream)
			r.write(data.Rd, data.Data)
			r.writeSig <- writeSignal{}
		}
	}()

	return nil
}

func (r *RegisterFile) read(addr uint8) datapath.XWord {
	if addr == 0 {
		return 0
	}
	r.m.RLock()
	defer r.m.RUnlock()
	return r.file[addr]
}

func (r *RegisterFile) write(addr uint8, data datapath.XWord) {
	if addr == 0 {
		return
	}

	r.m.Lock()
	defer r.m.Unlock()
	r.file[addr] = data
}

// Probe is a TEST-ONLY method that is used to read
// values from register address directly.
func (r *RegisterFile) Probe(addr uint8) datapath.XWord {
	return r.read(addr)
}

// Sideload is TEST-ONLY method used to load values directly into reg
func (r *RegisterFile) Sideload(addr uint8, val datapath.XWord) {
	r.write(addr, val)
}
