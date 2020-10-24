package cmd

import (
	"fmt"
	"log"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"
	"github.com/spf13/viper"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/bank"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api/transaction"
)

type ApiCmd struct {
	ynabClient ynab.ClientServicer
	path       string
	budgets    []Budget
}

func NewAPICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Push transactions to YNAB's api",
		RunE: func(cmd *cobra.Command, args []string) error {
			var yaml config
			viper.AddConfigPath(configPath)
			if err := viper.ReadInConfig(); err != nil {
				log.Fatal(err)
			}

			if err := viper.Unmarshal(&yaml); err != nil {
				log.Fatal(err)
			}

			dir, err := dirutil.DownloadPath()
			if err != nil {
				log.Fatal(err)
			}

			ApiCmd := ApiCmd{
				ynabClient: ynab.NewClient(yaml.Token),
				path:       dir,
				budgets:    yaml.Budgets,
			}

			return ApiCmd.run()
		},
	}
}

func (c ApiCmd) run() error {
	var responses []ResultResponse
	for _, budget := range c.budgets {
		var transactions []transaction.PayloadTransaction
		for _, account := range budget.Accounts {
			t, err := bank.ReadDir(c.path, account)
			if err != nil {
				return fmt.Errorf("failed extracting transactions: %w", err)
			}
			transactions = append(transactions, t...)
		}

		if len(transactions) < 1 {
			continue
		}

		createdTransactions, err := c.ynabClient.Transaction().CreateTransactions(budget.ID, transactions)
		if err != nil {
			log.Fatal(err)
		}

		responses = append(responses, ResultResponse{
			Budget:              budget,
			CreatedTransactions: createdTransactions,
		})
	}

	c.output(responses)

	return nil
}

func (c *ApiCmd) output(responses []ResultResponse) {
	t := tabby.New()
	t.AddHeader("Budget", "New", "Duplicated", "Total")
	for _, response := range responses {
		t.AddLine(response.Budget.Name, len(response.CreatedTransactions.TransactionIDs), len(response.CreatedTransactions.DuplicateImportIDs), len(response.CreatedTransactions.TransactionIDs)+len(response.CreatedTransactions.DuplicateImportIDs))
	}
	t.Print()
}

type Budget struct {
	ID       string
	Name     string
	Accounts []bank.Account
}

type ResultResponse struct {
	Budget
	*transaction.CreatedTransactions
}
