package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/csv"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api/transaction"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Push transactions to YNAB's api",
	Run: func(cmd *cobra.Command, args []string) {
		var ctx Context

		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		ctx.config.AccountID = os.Getenv("ACCOUNT_ID")
		ctx.config.AccesToken = os.Getenv("ACCES_TOKEN")
		ctx.config.BudgetID = os.Getenv("BUDGET_ID")
		ctx.config.IBAN = os.Getenv("IBAN")
		customPath, exists := os.LookupEnv("CUSTOM_PATH")
		if exists {
			ctx.config.CustomPath = &customPath
		}
		ctx.YNABClient = ynab.NewClient(ctx.config.AccesToken)
		createdTransactions, err := ctx.CSVToYNAB()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Transactions found: %d \n", len(createdTransactions.TransactionIDs)+len(createdTransactions.DuplicateImportIDs))
		log.Printf("Duplicated transactions: %d \n", len(createdTransactions.DuplicateImportIDs))
		log.Printf("-------------------------------- \n")
		log.Printf("Created transactions: %d \n", len(createdTransactions.TransactionIDs))
		os.Exit(0)
	},
}

type Context struct {
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

var errFailedToGetPath error = fmt.Errorf("failed to get path")

func (c Context) CSVToYNAB() (*transaction.CreatedTransactions, error) {
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

	exportLines, err := csv.GetLines(c.config.IBAN, *c.config.CustomPath)
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
