package testing

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// StreamFromFile loads a program in memory and surfaces a byte stream
func StreamFromFile(path string) (chan []byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open binary file: %w", err)
	}

	mem := bytes.NewBuffer(file)
	stream := make(chan []byte)

	go func() {
		defer close(stream)
		buf := make([]byte, 4)
		for {
			n, err := mem.Read(buf)
			if err == io.EOF {
				stream <- buf[:n]
				break
			}
			stream <- buf[:n]
		}
	}()

	return stream, nil
}
