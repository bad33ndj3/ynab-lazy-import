package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/joho/godotenv"

	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api/transaction"
)

type CSVToYNAB struct {
	YNABClient ynab.ClientServicer
	config
}

type config struct {
	AccountID  string
	AccesToken string
	BudgetID   string
	IBAN       string
	CustomPath *string
}

func main() {
	var cmd CSVToYNAB

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	cmd.config.AccountID = os.Getenv("ACCOUNT_ID")
	cmd.config.AccesToken = os.Getenv("ACCES_TOKEN")
	cmd.config.BudgetID = os.Getenv("BUDGET_ID")
	cmd.config.IBAN = os.Getenv("IBAN")
	customPath, exists := os.LookupEnv("CUSTOM_PATH")
	if exists {
		cmd.config.CustomPath = &customPath
	}
	cmd.YNABClient = ynab.NewClient(cmd.config.AccesToken)
	createdTransactions, err := cmd.CSVToYNAB()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transactions found: %d \n", len(createdTransactions.TransactionIDs)+len(createdTransactions.DuplicateImportIDs))
	log.Printf("Duplicated transactions: %d \n", len(createdTransactions.DuplicateImportIDs))
	log.Printf("-------------------------------- \n")
	log.Printf("Created transactions: %d \n", len(createdTransactions.TransactionIDs))
	os.Exit(0)
}

func (c CSVToYNAB) CSVToYNAB() (*transaction.CreatedTransactions, error) {
	if c.config.CustomPath == nil {
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("failed to get user for download path: %w", err)
		}
		downloadDir := fmt.Sprintf("%s/%s", usr.HomeDir, "Downloads")
		if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to get download path: %w", errFailedToGetPath)
		}

		c.config.CustomPath = &downloadDir
	}

	exportLines, err := getLines(c.config.IBAN, *c.config.CustomPath)
	if err != nil {
		return nil, err
	}

	var transactions []transaction.PayloadTransaction
	for _, line := range exportLines {
		trans, err := line.ToYNAB(c.config.AccountID)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *trans)
	}

	createdTransactions, err := c.YNABClient.Transaction().CreateTransactions(c.config.BudgetID, transactions)
	if err != nil {
		return nil, err
	}

	return createdTransactions, nil
}
