package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const configPath = "$HOME/.config/ynab-lazy-import"

type config struct {
	Token   string
	Budgets []budget
}

var rootCmd = &cobra.Command{
	Use:   "ynab-lazy-import",
	Short: "ynab-lazy-import is for lazy people who just want their csv to be uploaded to YNAB without lifting a finger!",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute executes the root Command, this will also register other commands in the init
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewAPICommand())
	initCmd, err := NewInitCommand()
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.AddCommand(initCmd)
}
