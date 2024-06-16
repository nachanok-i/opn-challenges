package utils

import (
	"encoding/csv"
	"fmt"
	"strings"
)

// ReadCSV reads CSV data from an io.Reader and returns the records
func ReadCSV(data []byte) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(string(data)))

	// Read all records
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV data: %w", err)
	}

	return records, nil
}
