package main

import (
	"fmt"
	"os"

	"github.com/nachanok-i/opn-challenges/utils"
)

func main() {
	// Get filename from command line argument
	if len(os.Args) < 2 {
		fmt.Println("Please enter file name.")
		return
	}

	fileName := os.Args[1]

	// Decode the file
	decodedData, err := utils.DecodeFile(fileName)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Process the decoded CSV data
	utils.ProcessFile(decodedData)
}
