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
func processTransaction(record *models.Tamboon, errorsChan chan error, report *models.Report, reportMu *sync.Mutex) {
	if limiter != nil {
		limiter.Wait(context.Background())
	}

	err := ChargeTransaction(record)
	if err != nil {
		if err.Error() == "(429/too_many_requests) API rate limit has been exceeded" {
			log.Debug("Rate limit exceeded")
			// Rate limit hit, apply rate limit and retry
			applyRateLimit()
			log.Debug("Rate limit hit, pausing for a while...")
			time.Sleep(2 * time.Second) // Adjust the sleep duration as necessary

			// Retry the transaction after pause
			err = ChargeTransaction(record)
			if err != nil {
				errorsChan <- fmt.Errorf("error retrying transaction for %v: %v", record.Name, err)
				reportMu.Lock()
				report.Failed += record.AmountSubunits
				reportMu.Unlock()
				return
			}
		} else {
			errorsChan <- fmt.Errorf("error charging transaction for %v: %v", record.Name, err)
			reportMu.Lock()
			report.Failed += record.AmountSubunits
			reportMu.Unlock()
			return
		}
	}
	reportMu.Lock()
	report.Success += record.AmountSubunits
	reportMu.Unlock()
}

// ProcessFile processes the decoded CSV data and calls the charge API
func ProcessFile(data []byte) *models.Report {
	log := GetLogger()
	fmt.Println("performing donations...")
	report := &models.Report{}
	var reportMu sync.Mutex

	// Number of worker goroutines
	var numWorkers int

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Debug("error loading .env file: %w", err)
		// Set default numWorkers to 10
		numWorkers = 10
	} else {
		numWorkers = atoi(os.Getenv("NUMBER_OF_WORKERS"))
	}

	// Create a CSV reader from decoded data
	reader := csv.NewReader(bytes.NewReader(data))

	// Read and ignore the header row
	if _, err := reader.Read(); err != nil {
		log.Debug("Error reading header: ", err)
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
			report.TotalReceived += transaction.AmountSubunits
			report.TotalDonator++
			if report.TopDonateAmount < transaction.AmountSubunits {
				report.TopDonateAmount = transaction.AmountSubunits
				report.TopDonors = []string{transaction.Name}
			} else if report.TopDonateAmount == transaction.AmountSubunits {
				report.TopDonors = append(report.TopDonors, transaction.Name)
			}
		}
	}()

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range transactionChan {
				log.Debug("record: ", record)
				processTransaction(record, errorsChan, report, &reportMu)
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
		log.Debug(err)
	}

	fmt.Println("done.")
	return report
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
