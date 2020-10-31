package cmd

import (
	"fmt"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"
	"github.com/spf13/viper"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/bank"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api/transaction"
)

// APICmd contains what the api command needs.
type APICmd struct {
	ynabClient ynab.ClientServicer
	path       string
	budgets    []budget
}

// NewAPICommand returns the API command as a cobra command.
func NewAPICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Push transactions to YNAB's api",
		RunE: func(cmd *cobra.Command, args []string) error {
			var yaml config
			path, err := dirutil.GetUserDirDirectory(configDirectory)
			if err != nil {
				return err
			}
			viper.AddConfigPath(path)
			if err := viper.ReadInConfig(); err != nil {
				return err
			}

			if err := viper.Unmarshal(&yaml); err != nil {
				return err
			}

			dir, err := dirutil.DownloadPath()
			if err != nil {
				return err
			}

			APICmd := APICmd{
				ynabClient: ynab.NewClient(yaml.Token),
				path:       dir,
				budgets:    yaml.Budgets,
			}

			return APICmd.run()
		},
	}
}

type budget struct {
	ID       string
	Name     string
	Accounts []bank.Account
}

type resultResponse struct {
	BudgetName            string
	NewTransactions       int
	DuplicateTransactions int
	TotalTransactions     int
}

func (c APICmd) run() error {
	responses := make([]resultResponse, len(c.budgets))
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
			return err
		}

		responses = append(responses, resultResponse{
			BudgetName:            budget.Name,
			NewTransactions:       len(createdTransactions.TransactionIDs),
			DuplicateTransactions: len(createdTransactions.DuplicateImportIDs),
			TotalTransactions:     len(createdTransactions.TransactionIDs) + len(createdTransactions.DuplicateImportIDs),
		})
	}

	c.output(responses)

	return nil
}

func (c *APICmd) output(responses []resultResponse) {
	t := tabby.New()
	t.AddHeader("budget", "New", "Duplicated", "Total")
	for _, response := range responses {
		t.AddLine(response.BudgetName, response.NewTransactions, response.DuplicateTransactions, response.TotalTransactions)
	}

	t.Print()
}
