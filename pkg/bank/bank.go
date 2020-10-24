package bank

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"
	"github.com/gocarina/gocsv"
	"go.bmvs.io/ynab/api/transaction"
)

type Bank interface {
	ToYNAB(accountID string) ([]transaction.PayloadTransaction, error)
	CorrectFile(path, iban string) bool
	Seperator() rune
}

type Account struct {
	Bank    string
	Name    string
	Account string
	Iban    string
}

func GetBank(bankType string) Bank {
	if strings.ToLower(bankType) == "ing" {
		return &INGLines{}
	}

	return &INGLines{}
}

func ReadDir(path string, account Account) ([]transaction.PayloadTransaction, error) {
	files, err := dirutil.FilePathWalkDir(path, ".csv")
	if err != nil {
		return nil, err
	}

	var transactions []transaction.PayloadTransaction
	for _, file := range files {
		lines, err := Read(file, account)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, lines...)
	}

	return transactions, nil
}

func Read(path string, account Account) ([]transaction.PayloadTransaction, error) {
	lines := GetBank(account.Bank)
	if !lines.CorrectFile(path, account.Iban) {
		return nil, nil
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = lines.Seperator()
		return r // Allows use pipe as delimiter
	})

	if err := gocsv.UnmarshalFile(file, lines); err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	transactions, err := lines.ToYNAB(account.Account)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
