package main

import (
	"flag"
	"fmt"
	"go.bmvs.io/ynab/api/transaction"
	"os"
	"os/user"
)

func main() {
	err := CSVToYNAB("test", flag.String("path", "", ""))
	if err != nil {
		panic(err)
	}
}

func CSVToYNAB(account string, dir *string) error {
	// get download dir
	if dir == nil {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get user for download path: %w", err)
		}
		downloadDir := fmt.Sprintf("%s/%s", usr.HomeDir, "Downloads")
		if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
			return errFailedToGetPath
		}

		dir = &downloadDir
	}

	// check for csv files
	exportLines, err := getLines(account, *dir)
	if err != nil {
		return err
	}

	// unmarshal exportfiles to ynab transactions
	var transactions []*transaction.PayloadTransaction
	for _, line := range exportLines {
		trans, err := line.ToYNAB("1")
		if err != nil {
			return err
		}
		transactions = append(transactions, trans)
	}

	// upload to ynab

	// optional delete uploaded files?

	return nil
}
