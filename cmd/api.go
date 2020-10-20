package cmd

import (
	"log"
	"os"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/csv"
	"github.com/bad33ndj3/ynab-lazy-import/pkg/downloaddirectory"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.bmvs.io/ynab"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Push transactions to YNAB's api",
	Run: func(cmd *cobra.Command, args []string) {
		var env config

		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		env.AccountID = os.Getenv("ACCOUNT_ID")
		env.AccesToken = os.Getenv("ACCES_TOKEN")
		env.BudgetID = os.Getenv("BUDGET_ID")
		env.IBAN = os.Getenv("IBAN")
		customPath, exists := os.LookupEnv("CUSTOM_PATH")
		if exists {
			env.CustomPath = &customPath
		}
		YNABClient := ynab.NewClient(env.AccesToken)
		if env.CustomPath == nil {
			dir, err := downloaddirectory.DownloadDirectory()
			if err != nil {
				log.Fatal(err)
			}
			env.CustomPath = dir
		}

		transactions, err := csv.GetLines(env.IBAN, *env.CustomPath, env.AccountID)
		if err != nil {
			log.Fatal(err)
		}

		createdTransactions, err := YNABClient.Transaction().CreateTransactions(env.BudgetID, transactions)
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

type config struct {
	AccountID  string
	AccesToken string
	BudgetID   string
	IBAN       string
	CustomPath *string
}
