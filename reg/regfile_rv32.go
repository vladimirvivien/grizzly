package reg

import (
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

type writeSignal = struct{}
type regfile = [datapath.RegSize]datapath.XWord
type RegisterFile struct {
	m sync.RWMutex
	file    *regfile
	writeSig chan writeSignal
	opInput <-chan datapath.OpFields
	dataInput <-chan datapath.RegisterData
	output  chan datapath.AluParam
}

func New() *RegisterFile {
	return &RegisterFile{
		writeSig: make(chan writeSignal),
		file: new(regfile),
		output: make(chan datapath.AluParam),
	}
}

func (r *RegisterFile) OpInput(in <-chan datapath.OpFields) {
	r.opInput = in
}

func (r *RegisterFile) DataInput(in <-chan datapath.RegisterData) {
	r.dataInput = in
}

func (r *RegisterFile) AluParamsOutput() <-chan datapath.AluParam {
	return r.output
}

func (r RegisterFile) Run() error {
	if r.opInput == nil {
		return fmt.Errorf("register file: missing opField input")
	}
	if r.dataInput == nil {
		return fmt.Errorf("register file: missing data input")
	}

	// Register instruction input loop
	// This loop handles instruction fields after decoding.
	// It setups control and prepare data for the ALU.
	// A semaphore signal is used to wait for writebacks
	// from operations (Op.R, Op.RI, Op.L, etc) that requires it.
	go func() {
		defer close(r.output)
		for op := range r.opInput {
			params := datapath.AluParam{
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
				r.output <- params
				<- r.writeSig

			case isa.Opcodes.RI:
				params.Op1 = r.read(op.Rs1)
				switch op.Funct3 {
				case integer.Slli.F3, integer.Srli.F3, integer.Srai.F3:
					params.Op2 = datapath.XWord(op.Shift)
				default:
					params.Op2 = op.Imm
				}

				r.output <- params
				<- r.writeSig
			}
		}
	}()

	// Register data input loop
	// This loop handles data that comes from any operation
	// that writes data back to the register (Opcode.R, Opcode.RI, Opcode.S, etc)
	// A semaphore signal is used to ensure that the read-loop can only happen
	// proceed after a previous write.
	go func(){
		for data := range r.dataInput  {
			r.write(data.Rd, data.Value)
			r.writeSig <- writeSignal{}
		}
	}()

	return nil
}

func (r *RegisterFile) read(addr uint8) datapath.XWord {
	r.m.Lock()
	defer r.m.Unlock()
	if addr == 0 {
		return 0
	}
	return r.file[addr]
}

func (r *RegisterFile) write(addr uint8, data datapath.XWord) {
	r.m.Lock()
	defer r.m.Unlock()
	if addr == 0 {
		return
	}
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
