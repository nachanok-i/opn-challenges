package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/nachanok-i/opn-challenges/cipher"
)

// DecodeFile reads and decodes the file content using Rot128Reader
func DecodeFile(fileName string) ([]byte, error) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer inputFile.Close()

	decoder, err := cipher.NewRot128Reader(inputFile)
	if err != nil {
		return nil, fmt.Errorf("error creating decoder: %w", err)
	}

	buf := make([]byte, 4096)
	var data []byte
	for {
		n, err := decoder.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading from decoder: %w", err)
		}
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
	}

	return data, nil
}
