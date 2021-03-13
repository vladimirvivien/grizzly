# Test Programs
This directory contains programs that can be used to test the core.

## RISC-V tool chain
As a pre-requisite, the RISC-V tool chain is needed to compile programs to the the RISCV ELF format.

* Download tool chain from GitHub - https://github.com/riscv/riscv-gnu-toolchain
* Alternatively, pre-built binaries can be download from Sifive - https://github.com/sifive/freedom-tools/releases

## Assembly programs
Test programs can be written in RISC-V assembly format.  Then the GNU tool chain can be used
to compile the programs to a binary file that can be executed by the core.

### Example
Assume the following is saved in source file `testprog.s`
```asm
main:
  addi x20, x0, 2
  addi x21, x0, 4
  addi x22, x0, 12
  addi x24, x20, 4
  add  x25, x21, x22
  slli x26, x21, 1
```

### RV32 Programs
Next, use the RISC-V GNU compiler to build a binary:

```shell
riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -march=rv32i -mabi=ilp32 -o add add.s
```

After that, use GNU's object-copy utility to generate a raw binary dump, from the compiled object file, which can be loaded in the simulated memory for execution.

```shell
riscv64-unknown-elf-objcopy --output-target=binary add add.bin
```

You can also get an assembly dump from the built object file using object dump utility to disassemble the built object file:
```shell
riscv64-unknown-elf-objdump -Ds add > add-dump.s
```

### RV64 Programs