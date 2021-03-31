main:
  addi x20, x0, 2
  addi x21, x0, 4
  addi x22, x0, 16
  add  x25, x20, x21
  sw   x25, 0(x22)
