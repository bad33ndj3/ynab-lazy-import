package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bmvs.io/ynab"
)

var rootCmd = &cobra.Command{
	Use:   "ynab-lazy-import",
	Short: "ynab-lazy-import is for lazy people who just want their csv to be uploaded to YNAB without lifting a finger!",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	var env config

	viper.AddConfigPath("$HOME/.config/ynab-lazy-import")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatal(err)
	}

	dir, err := dirutil.DownloadPath()
	if err != nil {
		log.Fatal(err)
	}

	ynabClient := ynab.NewClient(env.Token)
	rootCmd.AddCommand(NewAPICommand(ynabClient, dir, env.Budgets))
}
