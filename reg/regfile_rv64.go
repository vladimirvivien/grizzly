//go:build rv64 || rv64i

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
		InFields     datapath.Pin
		InAluData    datapath.Pin
		InMemData    datapath.Pin
		OutAluOps    datapath.Pin
		OutBranchOps datapath.Pin
	}{
		InFields:     datapath.Pin("regfile.in.opfields"),
		InAluData:    datapath.Pin("regfile.in.alu_data"),
		InMemData:    datapath.Pin("regfile.in.mem_data"),
		OutAluOps:    datapath.Pin("regfile.out.alu_ops"),
		OutBranchOps: datapath.Pin("regfile.out.branch_ops"),
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
	outBranch chan []byte
}

func New() *RegisterFile {
	reg := &RegisterFile{
		BaseComponent: datapath.NewBase(),
		writeSig:      make(chan writeSignal),
		file:          make(regfile, datapath.RegSize, datapath.RegSize),
		output:        make(chan []byte),
		outBranch:     make(chan []byte),
	}
	reg.Connect(Labels.OutAluOps, reg.output)
	reg.Connect(Labels.OutBranchOps, reg.outBranch)
	return reg
}

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

	go func() {
		defer close(r.output)
		for stream := range input {
			fields := datapath.DecodeOpFields(stream)
			op := datapath.Operation{
				Opcode: fields.Opcode,
				Rd:     fields.Rd,
			}

			switch fields.Opcode {
			case isa.Opcodes.R:
				op.AluOp = alu.EncodeAluOp(fields.Funct7, fields.Funct3)
				op.AluOperand1 = r.read(fields.Rs1)
				op.AluOperand2 = r.read(fields.Rs2)

				r.output <- datapath.EncodeOp(op)
				<-r.writeSig

			case isa.Opcodes.RI:
				op.AluOp = alu.EncodeAluOp(fields.Funct7, fields.Funct3)
				op.AluOperand1 = r.read(fields.Rs1)
				switch fields.Funct3 {
				case integer.Slli.F3, integer.Srli.F3, integer.Srai.F3:
					op.AluOperand2 = datapath.XWord(fields.Shift)
				default:
					op.AluOperand2 = datapath.XWord(fields.Imm)
				}

				r.output <- datapath.EncodeOp(op)
				<-r.writeSig

			case isa.Opcodes.L:
				op.AluOp = alu.Ops.Add
				op.AluOperand1 = r.read(fields.Rs1)
				op.AluOperand2 = datapath.XWord(fields.Imm)
				op.MemOp = fields.Funct3
				r.output <- datapath.EncodeOp(op)
				<-r.writeSig

			case isa.Opcodes.S:
				op.AluOp = alu.Ops.Add
				op.AluOperand1 = r.read(fields.Rs1)
				op.AluOperand2 = datapath.XWord(fields.Imm)
				op.MemOp = fields.Funct3
				op.MemData = r.read(fields.Rs2)
				r.output <- datapath.EncodeOp(op)

			case isa.Opcodes.J:
				op.AluOp = alu.Ops.Add
				op.PC = fields.PC
				op.AluOperand1 = fields.PC
				op.AluOperand2 = datapath.XWord(fields.Imm)
				r.output <- datapath.EncodeOp(op)
				<-r.writeSig

			case isa.Opcodes.JI:
				op.AluOp = alu.Ops.Add
				op.PC = fields.PC
				op.AluOperand1 = r.read(fields.Rs1)
				op.AluOperand2 = datapath.XWord(fields.Imm)
				r.output <- datapath.EncodeOp(op)
				<-r.writeSig

			case isa.Opcodes.B:
				brOp := datapath.BranchOp{
					PC:     fields.PC,
					Opcode: fields.Opcode,
					Funct3: fields.Funct3,
					RS1D:   r.read(fields.Rs1),
					RS2D:   r.read(fields.Rs2),
					Imm:    fields.Imm,
				}
				r.outBranch <- datapath.EncodeBranchOp(brOp)
			}
		}
	}()

	go func() {
		for dataStream := range inAluData {
			data := datapath.DecodeRegData(dataStream)
			r.write(data.Rd, data.Value)
			r.writeSig <- writeSignal{}
		}
	}()

	go func() {
		for dataStream := range inMemData {
			data := datapath.DecodeRegData(dataStream)
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

func (r *RegisterFile) Probe(addr uint8) datapath.XWord {
	return r.read(addr)
}

func (r *RegisterFile) Sideload(addr uint8, val datapath.XWord) {
	r.write(addr, val)
}
