package bank

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"
	"github.com/gocarina/gocsv"
	"go.bmvs.io/ynab/api/transaction"
)

// ErrNoValidDateString is thrown when a date string is not valid for YNAB.
var ErrNoValidDateString = fmt.Errorf("not a valid date string")

type Bank interface {
	ToYNAB(accountID string) ([]transaction.PayloadTransaction, error)
	CorrectFile(path, iban string) bool
	Separator() rune
}

type Account struct {
	Bank    string
	Name    string
	Account string
	Iban    string
}

var errNoValidBankType error = fmt.Errorf("failed getting valid bank type")

func GetBank(bankType string) (Bank, error) {
	if strings.EqualFold(bankType, "ing") {
		return &INGLines{}, nil
	}

	return &INGLines{}, errNoValidBankType
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
		for index := range lines {
			transactions = append(transactions, lines[index])
		}
	}

	return transactions, nil
}

func Read(path string, account Account) ([]transaction.PayloadTransaction, error) {
	lines, err := GetBank(account.Bank)
	if err != nil {
		return nil, err
	}
	if !lines.CorrectFile(path, account.Iban) {
		return nil, nil
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = lines.Separator()
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
