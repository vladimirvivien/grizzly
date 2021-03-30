package router

import (
	"fmt"

	"github.com/vladimirvivien/grizzly/datapath"
	"github.com/vladimirvivien/grizzly/isa"
)

var(
	Labels = struct {
		InAluResult datapath.Pin
		OutRegisterData datapath.Pin
	}{
		InAluResult: datapath.Pin("router.in.aluresult"),
		OutRegisterData: datapath.Pin("router.out.registerdata"),
	}
)

type Router struct {
	*datapath.BaseComponent
	regDataOutput chan []byte
}

func New() *Router {
	r := &Router{
		BaseComponent: datapath.NewBase(),
		regDataOutput: make(chan []byte),
	}
	r.Connect(Labels.OutRegisterData, r.regDataOutput)
	return r
}

func (r *Router) Run() error {
	aluResult := r.GetPin(Labels.InAluResult)
	if aluResult == nil {
		return fmt.Errorf("router: missing input: %s", Labels.InAluResult)
	}

	go func() {
		defer close(r.regDataOutput)
		for {
			input, opened := <-aluResult
			if !opened {
				return
			}

			result := datapath.DecodeResult(input)
			switch result.Opcode {
			case isa.Opcodes.R, isa.Opcodes.RI:
				// route to register file
				r.regDataOutput <- datapath.EncodeRegStore(datapath.RegisterStore{Rd: result.Rd, Data: result.AluOut})
			case isa.Opcodes.L:
			}

		}
	}()
	return nil
}
