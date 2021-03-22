package testing

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

// StreamFromFile loads a program in memory and surfaces a byte stream
func StreamFromFile(path string) (chan []byte, error) {
	mem, err := LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open binary file: %w", err)
	}

	stream := make(chan []byte)

	go func(f *bytes.Buffer, s chan []byte) {
		defer close(s)
		for {
			buf := make([]byte, 4) // 4 for XLEN=32
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			s <- buf[:n]
		}
	}(mem, stream)

	return stream, nil
}

func LoadFile(path string)(*bytes.Buffer, error){
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return &bytes.Buffer{}, fmt.Errorf("failed to load: %w", err)
	}

	return  bytes.NewBuffer(file), nil
}

func PrintRaw(buf *bytes.Buffer) {
	p := make([]byte, 4)
	for {
		n, err := buf.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Printf("%#v\n",p[:n])
	}
}

func PrintBinary(buf *bytes.Buffer) {
	p := make([]byte, 4)
	for {
		n, err := buf.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Printf("%032b\n",binary.LittleEndian.Uint32(p[:n]))
	}
}