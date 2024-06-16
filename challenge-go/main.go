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
	data, err := utils.DecodeFile(fileName)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Read the CSV data
	records, err := utils.ReadCSV(data)
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return
	}

	// Process the CSV records
	results, total, average, topDonors, err := utils.ProcessRecords(records)
	if err != nil {
		fmt.Println("Error processing CSV records:", err)
		return
	}

	// Output results
	fmt.Println("Total received: THB", total)
	fmt.Println("Average per person: THB", average)
	fmt.Println("Top donors:", topDonors)

	// Print the mapped structs and charge transactions
	for _, result := range results {
		chargeResp, err := utils.ChargeTransaction(&result)
		if err != nil {
			fmt.Println("Error charging transaction:", err)

			continue
		}
		fmt.Printf("Charge response: %+v\n", chargeResp)
	}
}
