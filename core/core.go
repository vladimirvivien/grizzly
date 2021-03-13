package core
//
//import (
//	"fmt"
//	"log"
//	"time"
//
//	"github.com/vladimirvivien/grizzly/alu"
//	"github.com/vladimirvivien/grizzly/clock"
//	"github.com/vladimirvivien/grizzly/ctrlunit"
//	"github.com/vladimirvivien/grizzly/device"
//	"github.com/vladimirvivien/grizzly/mem"
//	"github.com/vladimirvivien/grizzly/reg"
//)
//
//var (
//	In = struct {
//		Insts device.PinLabel
//	}{
//		Insts: "core.insts.in",
//	}
//
//	clk     = clock.New(10 * time.Millisecond)
//	memSize = uint64(4 * 1024)
//)
//
//type Core struct {
//	*device.Base
//	reg     device.Type
//	alu     device.Type
//	progMem device.Type
//	ctrl    device.ClockedType
//}
//
//func New() device.Type {
//	return newCore()
//}
//
//func newCore() *Core {
//	return &Core{
//		Base:    device.NewBase(),
//		reg:     reg.New(),
//		alu:     alu.New(),
//		progMem: mem.New(memSize),
//		ctrl:    ctrlunit.New(),
//	}
//}
//
//// Run starts the core and its components
//func (c *Core) Run() error {
//	log.Println("core: starting...")
//	if err := c.wireComponents(); err != nil {
//		return err
//	}
//	return c.startComponents()
//}
//
//func (c *Core) wireComponents() error {
//	log.Println("core: wiring components")
//	if c.GetPin(In.Insts) == nil {
//		return fmt.Errorf("instructions datapath not set")
//	}
//
//	// wire control unit
//	c.ctrl.SetPin(ctrlunit.In.Insts, c.GetPin(In.Insts))
//	c.ctrl.SetClock(clk)
//
//	// wire register
//	c.reg.SetPin(reg.In.RS1Addr, c.ctrl.GetPin(ctrlunit.Out.RS1))
//	c.reg.SetPin(reg.In.RS2Addr, c.ctrl.GetPin(ctrlunit.Out.RS2))
//	c.reg.SetPin(reg.In.Werf, c.ctrl.GetPin(ctrlunit.Out.Werf))
//	c.reg.SetPin(reg.In.RDAddr, c.ctrl.GetPin(ctrlunit.Out.RD))
//	rs2DataConnect := device.Connector("reg-alu-connector", c.reg.GetPin(reg.Out.RS2Data), 2)
//
//	// wire alu
//	c.alu.SetPin(alu.In.Operation, c.ctrl.GetPin(ctrlunit.Out.ALUOp))
//	c.alu.SetPin(alu.In.Operand1, c.reg.GetPin(reg.Out.RS1Data))
//	c.alu.SetPin(alu.In.Operand2, device.Mux(
//		"alu-op-mux",
//		c.ctrl.GetPin(ctrlunit.Out.ALUSrc),
//		rs2DataConnect[0],
//		c.ctrl.GetPin(ctrlunit.Out.Imm),
//	))
//	aluResultConnector := device.Connector("alu-mem-connector", c.alu.GetPin(alu.Out.Result), 2)
//
//	// memory
//	c.progMem.SetPin(mem.In.WriteEnable, c.ctrl.GetPin(ctrlunit.Out.MemWen))
//	c.progMem.SetPin(mem.In.ReadEnable, c.ctrl.GetPin(ctrlunit.Out.MemRen))
//	c.progMem.SetPin(mem.In.Operation, c.ctrl.GetPin(ctrlunit.Out.MemOp))
//	c.progMem.SetPin(mem.In.Address, aluResultConnector[0])
//	c.progMem.SetPin(mem.In.DataWrite, rs2DataConnect[1])
//
//	// alu-memory-mux to register writeback
//	aluMemWritebackMux := device.Mux(
//		"writeback-mux",
//		c.ctrl.GetPin(ctrlunit.Out.WBSel),
//		aluResultConnector[1],
//		c.progMem.GetPin(mem.Out.DataRead),
//	)
//
//	// wire alu-mem result writeback to reg
//	c.reg.SetPin(reg.In.Data, aluMemWritebackMux)
//
//	return nil
//}
//
//// startComponents loop through each component and invoke Run.
//func (c *Core) startComponents() error {
//	log.Println("core: starting components")
//	for _, comp := range []device.Type{c.ctrl, c.reg, c.alu, c.progMem} {
//		if err := comp.Run(); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
