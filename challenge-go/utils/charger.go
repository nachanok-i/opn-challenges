package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/nachanok-i/opn-challenges/models"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/sirupsen/logrus"
)

func ChargeTransaction(request *models.Tamboon) error {
	month, err := stringToMonth(request.ExpMonth)
	if err != nil {
		logrus.Debug(err)
	}
	yearInt, err := strconv.Atoi(request.ExpYear)
	if err != nil {
		return fmt.Errorf("invalid year: %s", request.ExpYear)
	}
	tokenRequest := &operations.CreateToken{
		Name:            request.Name,
		Number:          request.CCNumber,
		ExpirationMonth: time.Month(month),
		// I added 5 years in Expiration to make some transaction success
		ExpirationYear: yearInt + 5,
		SecurityCode:   request.CVV,
	}
	card, err := createTokenFromCard(tokenRequest)
	if err != nil {
		return err
	}
	chargeRequest := &operations.CreateCharge{
		Amount:   int64(request.AmountSubunits),
		Currency: "thb",
		Card:     card.ID,
	}
	_, err = chargeTransactionFromToken(chargeRequest)
	if err != nil {
		return err
	}

	return nil
}

// initializeOmiseClient initializes a new Omise client with API keys from environment variables
func initializeOmiseClient() (*omise.Client, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	publicKey := os.Getenv("OMISE_PUBLIC_KEY")
	secretKey := os.Getenv("OMISE_SECRET_KEY")

	if publicKey == "" || secretKey == "" {
		return nil, fmt.Errorf("OMISE_PUBLIC_KEY or OMISE_SECRET_KEY environment variable is not set")
	}

	client, err := omise.NewClient(publicKey, secretKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createTokenFromCard(request *operations.CreateToken) (*omise.Card, error) {
	client, err := initializeOmiseClient()
	if err != nil {
		return nil, err
	}

	result := &omise.Card{}
	err = client.Do(result, request)
	return result, err
}

func chargeTransactionFromToken(request *operations.CreateCharge) (*omise.Charge, error) {
	client, err := initializeOmiseClient()
	if err != nil {
		return nil, err
	}

	result := &omise.Charge{}
	err = client.Do(result, request)
	return result, err
}

func stringToMonth(monthStr string) (time.Month, error) {
	monthInt, err := strconv.Atoi(monthStr)
	if err != nil {
		return 0, fmt.Errorf("invalid month: %s", monthStr)
	}
	if monthInt < 1 || monthInt > 12 {
		return 0, fmt.Errorf("invalid month number: %d", monthInt)
	}
	return time.Month(monthInt), nil
}
