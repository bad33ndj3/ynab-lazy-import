package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api/transaction"
	"log"
	"os"
	"os/user"
)

type config struct {
	AccountID  string
	AccesToken string
	BudgetID   string
	IBAN       string
	CustomPath *string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var conf config
	conf.AccountID = os.Getenv("ACCOUNT_ID")
	conf.AccesToken = os.Getenv("ACCES_TOKEN")
	conf.BudgetID = os.Getenv("BUDGET_ID")
	conf.IBAN = os.Getenv("IBAN")
	customPath, exists := os.LookupEnv("CUSTOM_PATH")
	if exists {
		conf.CustomPath = &customPath
	}

	err = CSVToYNAB(conf)
	if err != nil {
		log.Fatal(err)
	}
}

func CSVToYNAB(conf config) error {
	// get download dir
	if conf.CustomPath == nil {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get user for download path: %w", err)
		}
		downloadDir := fmt.Sprintf("%s/%s", usr.HomeDir, "Downloads")
		if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
			return fmt.Errorf("failed to get download path: %w", errFailedToGetPath)
		}

		conf.CustomPath = &downloadDir
	}

	// check for csv files
	exportLines, err := getLines(conf.IBAN, *conf.CustomPath)
	if err != nil {
		return err
	}

	// unmarshal exportfiles to ynab transactions
	var transactions []transaction.PayloadTransaction
	for _, line := range exportLines {
		trans, err := line.ToYNAB(conf.AccountID)
		if err != nil {
			return err
		}
		transactions = append(transactions, *trans)
	}

	// upload to ynab
	client := ynab.NewClient(conf.AccesToken)
	createdTransactions, err := client.Transaction().CreateTransactions(conf.BudgetID, transactions)
	if err != nil {
		return err
	}

	log.Print(createdTransactions)
	// optional delete uploaded files?

	return nil
}
