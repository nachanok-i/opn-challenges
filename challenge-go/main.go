package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nachanok-i/opn-challenges/cipher"
)

func main() {
	// Get filename from command line argument
	if len(os.Args) < 2 {
		fmt.Println("Please enter file name.")
		return
	}

	fileName := os.Args[1]

	inputFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer inputFile.Close()

	decoder, err := cipher.NewRot128Reader(inputFile)
	if err != nil {
		fmt.Println("Error creating decoder: ", err)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err := decoder.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from decoder: ", err)
			return
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buf[:n]))
	}
}
