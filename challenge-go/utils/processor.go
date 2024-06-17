package utils

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/nachanok-i/opn-challenges/models"
	"golang.org/x/time/rate"
)

var limiter *rate.Limiter
var limiterMu sync.Mutex

// Function to apply rate limiting when necessary
func applyRateLimit() {
	limiterMu.Lock()
	defer limiterMu.Unlock()
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Limit(2), 5) // Adjust the rate limit as needed
	}
}

// Function to process a single transaction with retry logic
func processTransaction(record *models.Tamboon, errorsChan chan error) {
	if limiter != nil {
		limiter.Wait(context.Background())
	}

	err := ChargeTransaction(record)
	if err != nil {
		if err.Error() == "(429/too_many_requests) API rate limit has been exceeded" {
			fmt.Println("Rate limit exceeded")
			// Rate limit hit, apply rate limit and retry
			applyRateLimit()
			fmt.Println("Rate limit hit, pausing for a while...")
			time.Sleep(1 * time.Minute) // Adjust the sleep duration as necessary

			// Retry the transaction after pause
			err = ChargeTransaction(record)
			if err != nil {
				errorsChan <- fmt.Errorf("error retrying transaction for %v: %v", record.Name, err)
				return
			}
		} else {
			errorsChan <- fmt.Errorf("error charging transaction for %v: %v", record.Name, err)
			return
		}
	}
}

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

	// Read and ignore the header row
	if _, err := reader.Read(); err != nil {
		fmt.Printf("Error reading header: %v\n", err)
		return
	}

	// Create channels
	transactionChan := make(chan *models.Tamboon)
	errorsChan := make(chan error)
	wg := sync.WaitGroup{}

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
				fmt.Println("record: ", record)
				processTransaction(record, errorsChan)
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
