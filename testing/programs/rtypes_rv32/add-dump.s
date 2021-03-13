
add:     file format elf32-littleriscv

Contents of section .text:
 0000 130a2000 930a4000 130bc000 130c4a00  .. ...@.......J.
 0010 b38c6a01 139d1a00                    ..j.....        
Contents of section .riscv.attributes:
 0000 41190000 00726973 63760001 0f000000  A....riscv......
 0010 05727633 32693270 3000               .rv32i2p0.      

Disassembly of section .text:

00000000 <main>:
   0:	00200a13          	li	s4,2
   4:	00400a93          	li	s5,4
   8:	00c00b13          	li	s6,12
   c:	004a0c13          	addi	s8,s4,4
  10:	016a8cb3          	add	s9,s5,s6
  14:	001a9d13          	slli	s10,s5,0x1

Disassembly of section .riscv.attributes:

00000000 <.riscv.attributes>:
   0:	1941                	addi	s2,s2,-16
   2:	0000                	unimp
   4:	7200                	flw	fs0,32(a2)
   6:	7369                	lui	t1,0xffffa
   8:	01007663          	bgeu	zero,a6,14 <main+0x14>
   c:	0000000f          	fence	unknown,unknown
  10:	7205                	lui	tp,0xfffe1
  12:	3376                	fld	ft6,376(sp)
  14:	6932                	flw	fs2,12(sp)
  16:	7032                	flw	ft0,44(sp)
  18:	0030                	addi	a2,sp,8
