package utils

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/nachanok-i/opn-challenges/models"
	"golang.org/x/time/rate"
)

// Define the rate limit
const requestsPerSecond = 2 // Example: 2 requests per second
const burstLimit = 5        // Burst limit for rate limiter

// ProcessFile processes the decoded CSV data and calls the charge API
func ProcessFile(data []byte) {
	fmt.Println("performing donations...")

	// Number of worker goroutines
	var numWorkers int

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file: %w", err)
		// Set default numWorkers to 10
		numWorkers = 10
	} else {
		numWorkers = atoi(os.Getenv("NUMBER_OF_WORKERS"))
	}

	// Create a CSV reader from decoded data
	reader := csv.NewReader(bytes.NewReader(data))

	// Create channels
	transactionChan := make(chan *models.Tamboon)
	errorsChan := make(chan error)
	wg := sync.WaitGroup{}

	// Create a rate limiter
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burstLimit)

	// Start the reader goroutine
	go func() {
		defer close(transactionChan)
		for {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				errorsChan <- fmt.Errorf("error reading CSV: %v", err)
				return
			}
			// Convert record to models.Tamboon
			transaction := &models.Tamboon{
				Name:           record[0],
				AmountSubunits: atoi(record[1]),
				CCNumber:       record[2],
				CVV:            record[3],
				ExpMonth:       record[4],
				ExpYear:        record[5],
			}
			transactionChan <- transaction
		}
	}()

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range transactionChan {
				// Wait for permission from the rate limiter
				if err := limiter.Wait(context.Background()); err != nil {
					errorsChan <- fmt.Errorf("rate limiter error: %v", err)
					continue
				}

				// Process the record and call the charge API
				chargeResp, err := ChargeTransaction(record)
				if err != nil {
					errorsChan <- fmt.Errorf("error charging transaction for %v: %v", record.Name, err)
					continue
				}
				fmt.Printf("Charge response for %s: %+v\n", record.Name, chargeResp)
			}
		}()
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	// Handle errors
	for err := range errorsChan {
		fmt.Println(err)
	}

	fmt.Println("done.")
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
