package cmd

import (
	"log"
	"os"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/bank"
	"github.com/cheynewallace/tabby"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"

	"go.bmvs.io/ynab/api/transaction"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bmvs.io/ynab"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

type ResultResponse struct {
	Budget
	*transaction.CreatedTransactions
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Push transactions to YNAB's api",
	Run: func(cmd *cobra.Command, args []string) {
		var env config

		viper.AddConfigPath("$HOME/.config/ynab-lazy-import")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatal("no config file found at $HOME/.config/ynab-lazy-import")
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

		var responses []ResultResponse
		for _, budget := range env.Budgets {
			var transactions []transaction.PayloadTransaction
			for _, account := range budget.Accounts {
				t, err := bank.ReadDir(*env.CustomPath, account)
				if err != nil {
					log.Fatal(err)
				}
				transactions = append(transactions, t...)
			}

			if len(transactions) < 1 {
				return
			}
			createdTransactions, err := YNABClient.Transaction().CreateTransactions(budget.ID, transactions)
			if err != nil {
				log.Fatal(err)
			}

			responses = append(responses, ResultResponse{
				Budget:              budget,
				CreatedTransactions: createdTransactions,
			})
		}

		output(responses)
		os.Exit(0)
	},
}

func output(responses []ResultResponse) {
	t := tabby.New()
	t.AddHeader("Budget", "New", "Duplicated", "Total")
	for _, response := range responses {
		t.AddLine(response.Budget.Name, len(response.CreatedTransactions.TransactionIDs), len(response.CreatedTransactions.DuplicateImportIDs), len(response.CreatedTransactions.TransactionIDs)+len(response.CreatedTransactions.DuplicateImportIDs))
	}
	t.Print()
}

type config struct {
	Token      string
	Budgets    []Budget
	CustomPath *string
}

type Budget struct {
	ID       string
	Name     string
	Accounts []bank.Account
}
