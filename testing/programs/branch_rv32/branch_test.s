main:
  addi x1, x0, 5     # x1 = 5
  addi x2, x0, 5     # x2 = 5
  beq x1, x2, taken  # should branch to taken
  addi x3, x0, 1     # x3 = 1 (should not execute)
taken:
  addi x4, x0, 10    # x4 = 10
  addi x5, x0, 6     # x5 = 6
  bne x1, x5, taken2 # should branch to taken2 (5 != 6)
  addi x3, x0, 2     # x3 = 2 (should not execute)
taken2:
  addi x6, x0, 15    # x6 = 15
