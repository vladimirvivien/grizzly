package reg

import (
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/alu"
	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

var (
	Labels = struct {
		InFields  datapath.Pin
		InAluData datapath.Pin
		InMemData datapath.Pin
		OutAluOps datapath.Pin
	}{
		InFields:  datapath.Pin("regfile.in.opfields"),
		InAluData: datapath.Pin("regfile.in.alu_data"),
		InMemData: datapath.Pin("regfile.in.mem_data"),
		OutAluOps: datapath.Pin("regfile.out.alu_ops"),
	}
)

type writeSignal = struct{}
type regfile = []datapath.XWord
type RegisterFile struct {
	*datapath.BaseComponent
	m        sync.RWMutex
	file     regfile
	writeSig chan writeSignal
	output   chan []byte
}

func New() *RegisterFile {
	reg := &RegisterFile{
		BaseComponent: datapath.NewBase(),
		writeSig:      make(chan writeSignal),
		file:          make(regfile, datapath.RegSize, datapath.RegSize),
		output:        make(chan []byte),
	}
	reg.Connect(Labels.OutAluOps, reg.output)
	return reg
}

// Run starts the register file component
func (r *RegisterFile) Run() error {
	input := r.GetPin(Labels.InFields)
	if input == nil {
		return fmt.Errorf("register file: missing input: %s", Labels.InFields)
	}
	inAluData := r.GetPin(Labels.InAluData)
	if inAluData == nil {
		return fmt.Errorf("register file: missing data input: %s", Labels.InAluData)
	}
	inMemData := r.GetPin(Labels.InMemData)
	if inMemData == nil {
		return fmt.Errorf("register file: missing data input: %s", Labels.InMemData)
	}

	// Register instruction input loop
	// This input loop handles operation fields from the decoder.
	// A semaphore signal is used to wait for data_store
	// from operations (Op.R, Op.RI, Op.L, etc) which requires it.
	go func() {
		defer close(r.output)
		for stream := range input {
			fields := datapath.DecodeOpFields(stream)
			op := datapath.Operation{
				Opcode: fields.Opcode,
				Rd:     fields.Rd,
			}

			// Select ALU operands:
			switch fields.Opcode {
			case isa.Opcodes.R:
				op.AluOp = alu.EncodeAluOp(fields.Funct7, fields.Funct3)
				op.AluOperand1 = r.read(fields.Rs1)
				op.AluOperand2 = r.read(fields.Rs2)

				r.output <- datapath.EncodeOp(op)
				<-r.writeSig // wait for reg data writeback

			case isa.Opcodes.RI:
				op.AluOp = alu.EncodeAluOp(fields.Funct7, fields.Funct3)
				op.AluOperand1 = r.read(fields.Rs1)
				switch fields.Funct3 {
				case integer.Slli.F3, integer.Srli.F3, integer.Srai.F3:
					op.AluOperand2 = datapath.XWord(fields.Shift)
				default:
					op.AluOperand2 = fields.Imm
				}

				r.output <- datapath.EncodeOp(op)
				<-r.writeSig

			case isa.Opcodes.L:
				op.AluOp = alu.Ops.Add
				op.AluOperand1 = r.read(fields.Rs1)
				op.AluOperand2 = fields.Imm
				r.output <- datapath.EncodeOp(op)
				<-r.writeSig

			case isa.Opcodes.S:
				op.AluOp = alu.Ops.Add
				op.AluOperand1 = r.read(fields.Rs1)
				op.AluOperand2 = fields.Imm
				op.MemData = r.read(fields.Rs2)
				r.output <- datapath.EncodeOp(op)
			}
		}
	}()

	// ALU-to-register data store input loop
	// Handles register data storage writebacks from ALU for R-type operations.
	// A semaphore signal is used to ensure that the read-loop can only
	// proceed after a previous write.
	go func() {
		for dataStream := range inAluData {
			data := datapath.DecodeRegStore(dataStream)
			r.write(data.Rd, data.Value)
			r.writeSig <- writeSignal{}
		}
	}()

	// Memory-to-register data store input loop
	// Handles register data storage writebacks from Memory operations.
	// A semaphore signal is used to ensure that the read-loop can only
	// proceed after a previous write.
	go func() {
		for dataStream := range inMemData {
			data := datapath.DecodeRegStore(dataStream)
			r.write(data.Rd, data.Value)
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
