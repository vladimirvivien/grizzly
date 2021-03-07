package router

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

type Router struct {
	aluResultInput <-chan datapath.AluResult
	regDataOutput  chan datapath.RegisterData
}

func New() *Router {
	return &Router{
		regDataOutput: make(chan datapath.RegisterData),
	}
}

func (r *Router) AluResultInput(ch <-chan datapath.AluResult) {
	r.aluResultInput = ch
}

func (r *Router) RegisterDataOutput() <-chan datapath.RegisterData {
	return r.regDataOutput
}

func (r *Router) Run() error {
	if r.aluResultInput == nil {
		return fmt.Errorf("router: missing ALU result input")
	}

	go func() {
		defer close(r.regDataOutput)
		for {
			result, opened := <-r.aluResultInput
			if !opened {
				return
			}

			switch result.Opcode {
			case isa.Opcodes.R, isa.Opcodes.RI:
				// route to register file
				r.regDataOutput <- datapath.RegisterData{Rd: result.Rd, Value: result.Value}
			case isa.Opcodes.L:
			}

		}
	}()
	return nil
}
