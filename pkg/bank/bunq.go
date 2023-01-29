package bank

import (
	"fmt"
	"strconv"
	"strings"

	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
)

type BunqLines []bunq

func (i BunqLines) Separator() rune {
	return ','
}

func (i BunqLines) CorrectFile(path, iban string) bool {
	return strings.Contains(path, iban)
}

func (i BunqLines) ToYNAB(accountID string) ([]transaction.PayloadTransaction, error) {
	var lines []transaction.PayloadTransaction
	for _, line := range i {
		l, err := line.toYNAB(accountID)
		if err != nil {
			return nil, err
		}
		lines = append(lines, *l)
	}

	return lines, nil
}

type bunq struct {
	Datum            int    `csv:"Datum"`
	NaamOmschrijving string `csv:"Naam / Omschrijving"`
	Rekening         string `csv:"Rekening"`
	Tegenrekening    string `csv:"Tegenrekening"`
	Code             string `csv:"Code"`
	AfBij            string `csv:"Af Bij"`
	BedragEUR        string `csv:"Bedrag (EUR)"`
	Mutatiesoort     string `csv:"Mutatiesoort"`
	Mededelingen     string `csv:"Mededelingen"`
	SaldoNaMutatie   string `csv:"Saldo na mutatie"`
	Tag              string `csv:"Tag"`
}

func (e bunq) toYNAB(accountID string) (*transaction.PayloadTransaction, error) {
	trans := transaction.PayloadTransaction{
		AccountID: accountID,
		Cleared:   transaction.ClearingStatusCleared,
		Approved:  false,
		PayeeName: &e.NaamOmschrijving,
	}
	color := transaction.FlagColorGreen
	trans.FlagColor = &color
	if len(e.Mededelingen) > 195 {
		memo := e.Mededelingen[:195]
		trans.Memo = &memo
	} else {
		trans.Memo = &e.Mededelingen
	}
	amount, err := strconv.ParseInt(strings.ReplaceAll(e.BedragEUR, ",", ""), 10, 64)
	if err != nil {
		return nil, err
	}

	if e.AfBij == "Af" {
		amount = amount * -1
	}
	trans.Amount = amount * 10

	dateString := strconv.Itoa(e.Datum)
	if len(dateString) != 8 {
		return nil, fmt.Errorf("not a valid date string")
	}
	year := dateString[:4]
	month := dateString[4:6]
	day := dateString[6:]
	trans.Date, err = api.DateFromString(fmt.Sprintf("%s-%s-%s", year, month, day))
	if err != nil {
		return nil, err
	}

	importID := fmt.Sprintf("YNAB:%d:%s-%s-%s:1", amount, year, month, day)
	trans.ImportID = &importID

	return &trans, nil
}
