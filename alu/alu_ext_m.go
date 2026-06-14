//go:build ext_m

package alu

import (
	"github.com/vladimirvivien/grizzly/datapath"
)

func init() {
	// Register multiplication operations
	registerExtOp(Ops.Mul, func(op1, op2 datapath.XWord) datapath.XWord {
		return op1 * op2
	})
	registerExtOp(Ops.Mulh, mulh)
	registerExtOp(Ops.Mulhsu, mulhsu)
	registerExtOp(Ops.Mulhu, mulhu)

	// Register division operations
	registerExtOp(Ops.Div, func(op1, op2 datapath.XWord) datapath.XWord {
		o1 := int32(op1)
		o2 := int32(op2)
		if o2 == 0 {
			return 0xffffffff
		} else if o1 == -2147483648 && o2 == -1 {
			return uint32(o1)
		}
		return uint32(o1 / o2)
	})
	registerExtOp(Ops.Divu, func(op1, op2 datapath.XWord) datapath.XWord {
		if op2 == 0 {
			return 0xffffffff
		}
		return op1 / op2
	})
	registerExtOp(Ops.Rem, func(op1, op2 datapath.XWord) datapath.XWord {
		o1 := int32(op1)
		o2 := int32(op2)
		if o2 == 0 {
			return uint32(o1)
		} else if o1 == -2147483648 && o2 == -1 {
			return 0
		}
		return uint32(o1 % o2)
	})
	registerExtOp(Ops.Remu, func(op1, op2 datapath.XWord) datapath.XWord {
		if op2 == 0 {
			return op1
		}
		return op1 % op2
	})
}

// mulh** returns high 32-bit portion of multiplication product.
// For 32-bit operands, operation assumes 64-bit host machine.

// mulh interpret operands as signed
func mulh(data1, data2 datapath.XWord) datapath.XWord {
	result := (int64(data1) * int64(data2)) >> datapath.XWordLen
	return datapath.XWord(result)
}

// mulhsu interpret operands as signed/unsigned
func mulhsu(data1, data2 datapath.XWord) datapath.XWord {
	result := (uint64(int32(data1)) * uint64(data2)) >> datapath.XWordLen
	return datapath.XWord(result)
}

func mulhu(data1, data2 datapath.XWord) datapath.XWord {
	result := (uint64(data1) * uint64(data2)) >> datapath.XWordLen
	return datapath.XWord(result)
}
