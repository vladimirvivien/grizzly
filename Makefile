.PHONY: build-rv32i build-rv32im build-rv32imac build-rv64gc build-rv128i test clean

# 1. Baseline 32-bit (Microcontroller profile)
build-rv32i:
	go build -tags "rv32 rv32i" -o bin/grizzly main.go

# 2. 32-bit with Multiplication/Division (Grizzly's default)
build-rv32im:
	go build -tags "rv32 rv32i ext_m" -o bin/grizzly main.go

# 3. Embedded 32-bit with Atomics & Compressed
build-rv32imac:
	go build -tags "rv32 rv32i ext_m ext_a ext_c" -o bin/grizzly main.go

# 4. Linux-capable Application Server 64-bit profile (RV64GC)
build-rv64gc:
	go build -tags "rv64 rv64i ext_m ext_a ext_f ext_d ext_c" -o bin/grizzly main.go

# 5. Experimental 128-bit base profile
build-rv128i:
	go build -tags "rv128 rv128i" -o bin/grizzly main.go

# Run all tests using default RV32 configuration
test:
	go test ./...

clean:
	rm -rf bin/
