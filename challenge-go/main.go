package main

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
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
	report := utils.ProcessFile(decodedData)
	report.Average = float32(report.TotalReceived) / float32(report.TotalDonator)
	fmt.Printf("        total received: THB %15s\n", humanize.CommafWithDigits(float64(report.TotalReceived)/float64(100), 2))
	fmt.Printf("  successfully donated: THB %15s\n", humanize.CommafWithDigits(float64(report.Success)/float64(100), 2))
	fmt.Printf("       faulty donation: THB %15s\n", humanize.CommafWithDigits(float64(report.Failed)/float64(100), 2))
	fmt.Printf("\n")
	fmt.Printf("    average per person: THB %15s\n", humanize.CommafWithDigits(float64(report.Average), 2))
	fmt.Printf("            top donors: %s\n", report.TopDonors[0])
	for _, donator := range report.TopDonors[1:] {
		fmt.Printf("                        %s\n", donator)
	}
}
