package bank

import (
	"encoding/csv"
	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"
	"github.com/gocarina/gocsv"
	"go.bmvs.io/ynab/api/transaction"
	"io"
	"os"
	"strings"
)

type Bank interface {
	ToYNAB(accountID string) ([]transaction.PayloadTransaction, error)
	CorrectFile(path, iban string) bool
	Seperator() rune
}

type Account struct {
	Bank    string
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
	lines := GetBank(account.Bank)
	files, err := dirutil.FilePathWalkDir(path, ".csv")
	if err != nil {
		return nil, err
	}

	var transactions []transaction.PayloadTransaction
	for _, file := range files {
		lines, err := ReadFile(file, account.Iban, account.Account, lines)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, lines...)
	}

	return transactions, nil
}

func ReadFile(path, iban, accountID string, lines Bank) ([]transaction.PayloadTransaction, error) {
	err := Read(lines, path, iban)
	if err != nil {
		return nil, err
	}

	transactions, err := lines.ToYNAB(accountID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func Read(out Bank, path, iban string) error {
	if !out.CorrectFile(path, iban) {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = out.Seperator()
		return r // Allows use pipe as delimiter
	})

	if err := gocsv.UnmarshalFile(file, out); err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
