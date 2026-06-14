# Grizzly: A Clocked Concurrent RISC-V Simulator

Grizzly is a concurrent, clocked RISC-V CPU simulator written in Go. Unlike serial simulators that execute instructions in a single thread, Grizzly models individual CPU components (Program Counter, Instruction Memory, Decoder, Register File, ALU, Brancher, Data Memory) as independent execution stages running concurrently in their own goroutines. 

These components communicate asynchronously over Go channels that act as hardware "wires" carrying serialized datapath packets. A central clock component synchronizes state transitions across stage boundaries.

---

## 1. Architecture and Core Design

Grizzly is designed around a concurrent, message-passing pipeline that models hardware signals as byte-stream packets.

### Component Abstraction
Every CPU component implements the `datapath.Component` interface:

```go
type Component interface {
    Connect(label Pin, stream Bytestream)
    GetPin(label Pin) Bytestream
    Run() error
}
```

Components interact through three primary abstractions:
- **`datapath.Pin`**: A string identifier representing an entry or exit connection port (e.g. `pc.in.pc_operation`, `alu.out.reg_data`).
- **`datapath.Bytestream`**: A read-only byte channel (`<-chan []byte`) modeling a physical bus or pin.
- **`datapath.BaseComponent`**: A helper struct providing thread-safe storage for pin mappings.

### Hardware Communication & Serialization
Components do not share state or pass pointers. Instead, data is serialized into byte streams before transmission and deserialized by the receiving component. This enforces component isolation and allows the pipeline to scale transparently across different processor widths (32-bit, 64-bit, and 128-bit). 

For example, when the ALU outputs register writeback data, the values are encoded into a byte slice using the `encoding/binary` package and placed onto the `alu.out.reg_data` pin:

```go
// Encoding register writeback data
buf := make([]byte, 1 + XWordBytes)
buf[0] = rd
binary.LittleEndian.PutUint64(buf[1:], value)
```

The Register File receives this stream, decodes the register index and value, and commits the state.

### Clock Synchronization
Although stage execution is asynchronous, pipeline progression is synchronized by a simulation clock (`clock.Clock`).
- The clock component emits ticks on a channel at a configured interval.
- Clocked boundary components (such as the Program Counter) synchronize state updates with clock ticks.
- The Program Counter blocks until a tick is received, fetches the resolved program address, and sends it downstream, triggering execution stages in step.

### Multiplexing (MUX)
Physical hardware multiplexers are implemented using asynchronous channel-merging loops. For example, register operations and brancher operations are multiplexed into the ALU using a merge utility that funnels multiple input channels into a single output stream:

```go
func merge(ch1, ch2 datapath.Bytestream) datapath.Bytestream {
    out := make(chan []byte)
    go func() {
        defer close(out)
        var wg sync.WaitGroup
        wg.Add(2)
        go func() { defer wg.Done(); for v := range ch1 { out <- v } }()
        go func() { defer wg.Done(); for v := range ch2 { out <- v } }()
        wg.Wait()
    }()
    return out
}
```

### Component Interconnection Diagram
The following diagram shows the connections and pin names between the simulation stages:

```
                                 +--------------+
                                 |      PC      |<──────────────────────────────┐
                                 +-------+------+                               │
                                         | (pc.out.counter)                     │
                                         v                                      │
                                 +-------+------+                               │
                                 |  Inst Memory |                               │
                                 +-------+------+                               │
                                         | (mem.out.instruction)                │
                                         v                                      │
                                 +-------+------+                               │
                                 |    Decoder   |                               │
                                 +-------+------+                               │
                                         | (decoder.out.fields)                 │
                                         v                                      │
                                 +-------+------+                               │
                  +─────────────>| RegisterFile |<─────────────────────────┐    │
                  │              +───┬──────┬───+                          │    │
                  │                  │      │                              │    │
                  │ (reg.out.alu_ops)│      │(reg.out.branch_ops)          │    │
                  │                  │      v                              │    │
                  │                  │  +───┴───+                          │    │
                  │                  │  |Branchr|                          │    │
                  │                  │  +───┬───+                          │    │
                  │                  │      │ (branch.out.operation)       │    │
                  │                  v      v                              │    │
                  │                +─┴──────┴─+                            │    │
                  │                |  (merge) |                            │    │
                  │                +────┬─────+                            │    │
                  │                     | (alu.in.operations)              │    │
                  │                     v                                  │    │
                  │                 +───┴───+                              │    │
                  │                 |  ALU  |                              │    │
                  │                 +─┬──┬──+                              │    │
                  │                   │  │ (alu.out.mem_op)                │    │
                  │                   │  v                                 │    │
                  │  (alu.out.reg_d)  │ +┴──────────+                      │    │
                  │                   │ | Data Mem  |                      │    │
                  │                   │ +────┬──────+                      │    │
                  │                   │      │ (mem.out.reg_data)          │    │
                  │                   │      └─────────────────────────────┘    │
                  │                   │                                         │
                  │                   │ (alu.out.pc_op)                         │
                  │                   └─────────────────────────────────────────┘
                  │
                  +─────────────────────────────────────────────────────────────┘
```

---

## 2. Modular Design & Build Tags

Grizzly implements different RISC-V architectural profiles and word widths using Go build constraints (`//go:build`). This separates target-specific definitions from the core pipeline logic.

### Layer 1: Width-Class (`rv32`, `rv64`, `rv128`)
Define basic data types and sizes for registers and memory addresses:
- **`datapath/types_rv32.go`**: Defines `XWord` as `uint32` and `RegSize` as `32`.
- **`datapath/types_rv64.go`**: Defines `XWord` as `uint64` and `RegSize` as `64`.
- **Memory Alignment**: Swapped dynamically: `base_rv32.go` asserts 4-byte boundaries, whereas `base_rv64.go` asserts 8-byte boundaries for doubleword operations.

### Layer 2: Instruction-Set Base (`rv32i`, `rv32e`, `rv64i`)
Determine instruction decoding rules and register file capacities (e.g. `rv32e` allocates 16 registers; `rv32i` allocates 32 registers).

### Layer 3: Feature Extensions (`ext_m`, `ext_a`, `ext_f`, `ext_d`, `ext_c`, `ext_v`)
Features like Multiplication (`M` extension) or Atomics (`A` extension) are isolated in separate files.
- The ALU execution function (`aluFunc`) exposes a registry hook for optional operations:
  ```go
  var extOperations = make(map[uint8]func(op1, op2 datapath.XWord) datapath.XWord)
  ```
- Optional files (e.g. `alu_ext_m.go` under `//go:build ext_m`) use Go's `init()` function to register their execution blocks.
- Compiling without the build tag omits the extension files from compilation entirely, generating a smaller binary without multiplication/division instructions.

---

## 3. Building the Simulator

The compilation profile is controlled by passing the target build tags using the `go build` command. A root `Makefile` is provided to manage these profiles.

### Compilation Commands

To build specific RISC-V configurations, execute the corresponding Makefile target:

```bash
# Compile baseline 32-bit integer simulator
make build-rv32i

# Compile 32-bit with Multiplication/Division extension (Grizzly's default)
make build-rv32im

# Compile embedded 32-bit with Atomics and Compressed instructions
make build-rv32imac

# Compile Linux-capable 64-bit application server profile (RV64GC)
make build-rv64gc

# Compile experimental 128-bit base profile
make build-rv128i
```

Alternatively, invoke `go build` directly with the desired tags:
```bash
go build -tags "rv64 rv64i ext_m ext_a" -o bin/grizzly main.go
```

To clean build output:
```bash
make clean
```

---

## 4. Compiling and Running RISC-V Programs

Grizzly executes programs stored as raw, freestanding RISC-V binary payloads. 

### Cross-Compiling Programs with Zig
You can write tests in C or assembly and cross-compile them using the local `zig` compiler.

#### 1. Target: 32-Bit RISC-V (`rv32i`)
Create an assembly source file `program.s`:
```assembly
main:
  addi x10, x0, 5
  addi x11, x0, 10
  add  x12, x10, x11
```

Cross-compile and extract the binary payload:
```bash
# Compile to a freestanding ELF binary
zig cc -target riscv32-freestanding-none -mcpu=generic_rv32 -nostdlib -o program.elf program.s

# Extract raw instruction bytes
zig objcopy -O binary program.elf program.bin
```

#### 2. Target: 64-Bit RISC-V (`rv64i`)
Cross-compile and extract the binary payload:
```bash
# Compile to a freestanding ELF binary
zig cc -target riscv64-freestanding-none -mcpu=generic_rv64 -nostdlib -o program.elf program.s

# Extract raw instruction bytes
zig objcopy -O binary program.elf program.bin
```

### Loading and Executing Binaries
To execute a program, load the generated `.bin` payload into the simulator's instruction memory:

```go
// Read binary payload file
content, err := ioutil.ReadFile("program.bin")
if err != nil {
    log.Fatal(err)
}

// Instantiate core and load program starting at address 0
cpu := core.New()
cpu.imem.SetStore(content)

// Start pipeline execution
if err := cpu.Run(); err != nil {
    log.Fatal(err)
}
```

### Simulator Termination
When the simulator reaches the end of the loaded binary payload, the Program Counter continues to tick. To prevent out-of-bounds panics, Grizzly's instruction memory returns a standard RISC-V NOP (`0x00000013`, `addi x0, x0, 0`) for any addresses beyond the program size. This keeps the pipeline running safely, allowing developers to inspect register and memory states.

---

## 5. Testing and Verification

To run unit and integration tests:

### Run 32-bit test suite
```bash
make test
# or:
go test -tags "rv32 rv32i" ./...
```

### Run 64-bit test suite
```bash
go test -tags "rv64 rv64i" ./...
```
