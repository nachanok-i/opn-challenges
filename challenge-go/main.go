package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/nachanok-i/opn-challenges/cipher"
)

type Tamboon struct {
	Name           string
	AmountSubunits int
	CCNumber       string
	CVV            string
	ExpMonth       string
	ExpYear        string
}

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
	var data []byte
	for {
		n, err := decoder.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from decoder: ", err)
			return
		}
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
		// fmt.Print(string(buf[:n]))
	}

	r := csv.NewReader(strings.NewReader(string(data)))

	// Read CSV headers
	headers, err := r.Read()
	if err != nil {
		fmt.Println("Error reading CSV headers:", err)
		return
	}
	fmt.Println("Headers:", headers)

	var total, top int
	// success, failed,
	var average float64
	var topDonors []string
	// Read CSV records and map them to the struct
	var results []Tamboon
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV record:", err)
			return
		}

		amount, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println("Error converting AmountSubunits to int:", err)
			return
		}

		result := Tamboon{
			Name:           record[0],
			AmountSubunits: amount,
			CCNumber:       record[2],
			CVV:            record[3],
			ExpMonth:       record[4],
			ExpYear:        record[5],
		}
		results = append(results, result)
		total += amount
		if top < amount {
			top = amount
			// clear old topDonors
			topDonors = []string{}
			topDonors = append(topDonors, result.Name)
		} else if top == amount {
			topDonors = append(topDonors, result.Name)
		}
	}

	average = float64(total) / float64(len(results))

	fmt.Println("total received: THB ", total)
	fmt.Println("average per person: THB ", average)
}
