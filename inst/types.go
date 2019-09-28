package inst

type Type uint32

type Fields struct {
	Opcode uint32 // opcode format
	Rd     uint32 // destination register
	Funct3 uint32 // branch func
	Rs1    uint32 // register 1
	Rs2    uint32 // register 2
	Funct7 uint32 // alu func
	Imm    uint32 // immediate
}

type rop struct {
	Name string
	Funct7,
	Funct3,
	Opcode uint32
}

var (
	Opcodes = struct {
		R,
		L,
		RI,
		Ecall,
		S,
		SB,
		U,
		UJ uint32
	}{
		R:     0b0110011,
		L:     0x03,
		RI:    0x13,
		Ecall: 0x73,
		S:     0x23,
		SB:    0x63,
		U:     0x37,
		UJ:    0x6F,
	}

	Add  = rop{Name: "add", Funct7: 0b0000000, Funct3: 0b000, Opcode: Opcodes.R}
	Sub  = rop{Name: "sub", Funct7: 0b0100000, Funct3: 0b000, Opcode: Opcodes.R}
	Sll  = rop{Name: "sll", Funct7: 0b0000000, Funct3: 0b001, Opcode: Opcodes.R}
	Slt  = rop{Name: "slt", Funct7: 0b0000000, Funct3: 0b010, Opcode: Opcodes.R}
	Sltu = rop{Name: "sltu", Funct7: 0b0000000, Funct3: 0b011, Opcode: Opcodes.R}
	Xor  = rop{Name: "xor", Funct7: 0b0000000, Funct3: 0b100, Opcode: Opcodes.R}
	Srl  = rop{Name: "slr", Funct7: 0b0000000, Funct3: 0b101, Opcode: Opcodes.R}
	Sra  = rop{Name: "sra", Funct7: 0b0100000, Funct3: 0b101, Opcode: Opcodes.R}
	Or   = rop{Name: "or", Funct7: 0b0000000, Funct3: 0b110, Opcode: Opcodes.R}
	And  = rop{Name: "and", Funct7: 0b0000000, Funct3: 0b111, Opcode: Opcodes.R}
)
