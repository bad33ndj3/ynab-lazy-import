package csv

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/bad33ndj3/ynab-lazy-import/pkg/bank"
	"github.com/bad33ndj3/ynab-lazy-import/pkg/dirutil"

	"github.com/gocarina/gocsv"
	"go.bmvs.io/ynab/api/transaction"
)

func ReadDir(path, iban, accountID string) ([]transaction.PayloadTransaction, error) {
	files, err := dirutil.FilePathWalkDir(path, ".csv")
	if err != nil {
		return nil, err
	}

	var transactions []transaction.PayloadTransaction
	for _, file := range files {
		lines, err := ReadFile(file, iban, accountID)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, lines...)
	}

	return transactions, nil
}

func ReadFile(path, iban, accountID string) ([]transaction.PayloadTransaction, error) {
	var transactions []transaction.PayloadTransaction
	if !strings.Contains(path, iban) {
		return nil, nil
	}

	lines := make([]*bank.INGExport, 1)
	err := Read(path, ';', &lines)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		trans, err := line.ToYNAB(accountID)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *trans)
	}

	return transactions, nil
}

func Read(path string, seperator rune, out interface{}) error {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = seperator
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
