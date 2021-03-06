package reg

import (
	"fmt"
	"sync"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
	"github.com/vladimirvivien/grizzly/isa/integer"
)

type RegisterFile struct {
	sync.RWMutex
	file    []datapath.XWord
	opInput <-chan datapath.OpFields
	dataInput <-chan datapath.RegisterData
	output  chan datapath.AluParam
}

func New() *RegisterFile {
	return &RegisterFile{
		file:   make([]datapath.XWord, 32, 32),
		output: make(chan datapath.AluParam),
	}
}

func (r *RegisterFile) OpInput(in <-chan datapath.OpFields) {
	r.opInput = in
}

func (r *RegisterFile) DataInput(in <-chan datapath.RegisterData) {
	r.dataInput = in
}

func (r *RegisterFile) AluParams() <-chan datapath.AluParam {
	return r.output
}

func (r RegisterFile) Run() error {
	if r.opInput == nil {
		return fmt.Errorf("register file missing opField input")
	}
	if r.dataInput == nil {
		return fmt.Errorf("register file missing data input")
	}

	// operation fields opInput loop
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
			case isa.Opcodes.RI:
				params.Op1 = r.read(op.Rs1)
				switch op.Funct3 {
				case integer.Slli.F3, integer.Srli.F3, integer.Srai.F3:
					params.Op2 = datapath.XWord(op.Shift)
				default:
					params.Op2 = op.Imm
				}
			}

			r.output <- params
		}
	}()

	// register data input
	go func(){
		for data := range r.dataInput  {
			r.write(data.Rd, data.Value)
		}
	}()

	return nil
}

func (r *RegisterFile) read(addr uint8) (data datapath.XWord) {
	r.RLock()
	defer r.RUnlock()
	if addr == 0 {
		data = 0
	} else {
		data = r.file[addr]
	}
	return
}

func (r *RegisterFile) write(addr uint8, data datapath.XWord) {
	r.Lock()
	defer r.Unlock()
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
