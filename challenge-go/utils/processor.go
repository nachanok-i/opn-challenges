package utils

import (
	"fmt"
	"strconv"

	"github.com/nachanok-i/opn-challenges/models"
)

// ProcessRecords processes CSV records and maps them to Tamboon structs
func ProcessRecords(records [][]string) ([]models.Tamboon, int, float64, []string, error) {
	if len(records) < 2 {
		return nil, 0, 0, nil, fmt.Errorf("no data in CSV records")
	}

	headers := records[0]
	fmt.Println("Headers:", headers)

	var results []models.Tamboon
	var total, top int
	var topDonors []string

	for _, record := range records[1:] {
		amount, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, 0, 0, nil, fmt.Errorf("error converting AmountSubunits to int: %w", err)
		}

		result := models.Tamboon{
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
			topDonors = []string{result.Name}
		} else if top == amount {
			topDonors = append(topDonors, result.Name)
		}
	}

	average := float64(total) / float64(len(results))
	return results, total, average, topDonors, nil
}
