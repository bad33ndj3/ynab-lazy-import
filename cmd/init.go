package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/bank"
	"github.com/spf13/cobra"
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api/account"
	"gopkg.in/yaml.v2"
)

// InitCmd contains what the Init command needs
type InitCmd struct {
	ynab ynab.ClientServicer
}

var errFileAlreadyExists = fmt.Errorf("config file already exists")

// NewInitCommand creates the config file
func NewInitCommand() (*cobra.Command, error) {
	var token string
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create a config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(configPath); err == nil {
				return errFileAlreadyExists
			}

			ApiCmd := InitCmd{
				ynab: ynab.NewClient(token),
			}

			return ApiCmd.run(token)
		},
	}
	initCmd.Flags().StringVarP(&token, "token", "t", "", "YNAB Token (required)")
	if err := initCmd.MarkFlagRequired("token"); err != nil {
		return nil, err
	}

	return initCmd, nil
}

func (c InitCmd) run(token string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	path := strings.Replace(configPath, "$HOME", usr.HomeDir, 1) + "/config.yaml"
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config file already exists")
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
	} else {
		return err
	}

	budgetList, err := c.ynab.Budget().GetBudgets()
	if err != nil {
		return err
	}

	conf := config{
		Token: token,
	}
	for _, b := range budgetList {
		accountList, err := c.ynab.Account().GetAccounts(b.ID)
		if err != nil {
			return err
		}
		var accounts []bank.Account
		for _, a := range accountList {
			if a.Type != account.TypeChecking {
				continue
			}
			accounts = append(accounts, bank.Account{
				Account: a.ID,
				Name:    a.Name,
			})
		}

		conf.Budgets = append(conf.Budgets, budget{
			ID:       b.ID,
			Name:     b.Name,
			Accounts: accounts,
		})
	}

	out, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write(out)
	if err != nil {
		return err
	}

	return nil
}
