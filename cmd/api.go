package cmd

import (
	"log"
	"os"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"

	"go.bmvs.io/ynab/api/transaction"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/csv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		viper.AddConfigPath("$HOME/.config/ynab-lazy-import")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatal("no config file found at $HOME/.ynab")
			} else {
				log.Fatal(err)
			}
		}

		err := viper.Unmarshal(&env)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}

		YNABClient := ynab.NewClient(env.Token)

		if env.CustomPath == nil {
			dir, err := dirutil.DownloadPath()
			if err != nil {
				log.Fatal(err)
			}
			env.CustomPath = dir
		}

		for _, budget := range env.Budgets {
			var transactions []transaction.PayloadTransaction
			for _, account := range budget.Accounts {
				t, err := csv.ReadDir(*env.CustomPath, account.Iban, account.Account)
				if err != nil {
					log.Fatal(err)
				}
				transactions = append(transactions, t...)
			}

			if len(transactions) < 1 {
				return
			}
			createdTransactions, err := YNABClient.Transaction().CreateTransactions(budget.Budget, transactions)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("-------------------------------- \n")
			log.Printf("Transactions found: %d \n", len(createdTransactions.TransactionIDs)+len(createdTransactions.DuplicateImportIDs))
			log.Printf("Duplicated transactions: %d \n", len(createdTransactions.DuplicateImportIDs))
			log.Printf("Created transactions: %d \n", len(createdTransactions.TransactionIDs))
			log.Printf("\n")
		}

		os.Exit(0)
	},
}

type config struct {
	Token   string
	Budgets []struct {
		Budget   string
		Accounts []struct {
			Account string
			Iban    string
		}
	}

	CustomPath *string
}
