package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
)

type CSVTestSuite struct {
	suite.Suite

	testAssets string
}

func (t *CSVTestSuite) SetupSuite() {
	t.testAssets = "testassets"
}

func (t *CSVTestSuite) TestCSVToINGExport() {
	tests := []struct {
		expectedOutput []*INGExport
		inAccount      string
		inDir          string
		err            error
	}{
		{
			expectedOutput: []*INGExport{
				{
					Datum:            20200101,
					NaamOmschrijving: "Origin.com EA",
					Rekening:         "NL13INGB0000000000",
					Tegenrekening:    "NL14RABO0000000000",
					Code:             "ID",
					AfBij:            "Af",
					BedragEUR:        "3,99",
					Mutatiesoort:     "iDEAL",
					Mededelingen:     "example",
					SaldoNaMutatie:   "20,00",
					Tag:              "",
				},
			},
			inAccount: "NL13INGB0000000000",
			inDir:     fmt.Sprintf("%s/%s", t.testAssets, "base"),
		},
		{
			inAccount: "accountNonExisting",
			inDir:     fmt.Sprintf("%s/%s", t.testAssets, "base"),
		},
		{
			inAccount: "accountNonExisting",
			inDir:     fmt.Sprintf("%s/%s", t.testAssets, "pathNonExisting"),
			err:       errFailedToGetPath,
		},
	}
	for _, test := range tests {
		t.Run("", func() {
			lines, err := getLines(test.inAccount, test.inDir)
			t.Require().Equal(test.err, err)
			t.Require().Equal(test.expectedOutput, lines)
		})
	}
}

func (t *CSVTestSuite) TestToYNAB() {
	tests := []struct {
		Line      INGExport
		AccountID string
		PayeeName string
		Memo      string
		FlagColor transaction.FlagColor
		Date      time.Time
		ImportID  string
		Amount    int64
	}{
		{
			Line: INGExport{
				Datum:            20200102,
				NaamOmschrijving: "Shopping",
				Rekening:         "NL13INGB0000000000",
				Tegenrekening:    "NL14RABO0000000000",
				Code:             "ID",
				AfBij:            "Af",
				BedragEUR:        "3,99",
				Mutatiesoort:     "iDEAL",
				Mededelingen:     "example",
				SaldoNaMutatie:   "20,00",
			},
			AccountID: "test_accountid",
			PayeeName: "Shopping",
			Memo:      "example",
			FlagColor: transaction.FlagColorGreen,
			Date:      time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			ImportID:  "YNAB:-399:2020-01-02:1",
			Amount:    -3990,
		},
		{
			Line: INGExport{
				Datum:            20200102,
				NaamOmschrijving: "Shopping",
				Rekening:         "NL13INGB0000000000",
				Tegenrekening:    "NL14RABO0000000000",
				Code:             "ID",
				AfBij:            "Af",
				BedragEUR:        "123,99",
				Mutatiesoort:     "iDEAL",
				Mededelingen:     "example",
				SaldoNaMutatie:   "20,00",
			},
			AccountID: "test_accountid",
			PayeeName: "Shopping",
			Memo:      "example",
			FlagColor: transaction.FlagColorGreen,
			Date:      time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			ImportID:  "YNAB:-12399:2020-01-02:1",
			Amount:    -123990,
		},
		{
			Line: INGExport{
				Datum:            20200102,
				NaamOmschrijving: "Shopping",
				Rekening:         "NL13INGB0000000000",
				Tegenrekening:    "NL14RABO0000000000",
				Code:             "ID",
				AfBij:            "Bij",
				BedragEUR:        "22,00",
				Mutatiesoort:     "iDEAL",
				Mededelingen:     "example",
				SaldoNaMutatie:   "20,00",
			},
			AccountID: "test_accountid",
			PayeeName: "Shopping",
			Memo:      "example",
			FlagColor: transaction.FlagColorGreen,
			Date:      time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			ImportID:  "YNAB:2200:2020-01-02:1",
			Amount:    22000,
		},
		{
			Line: INGExport{
				Datum:            20200102,
				NaamOmschrijving: "Shopping",
				Rekening:         "NL13INGB0000000000",
				Tegenrekening:    "NL14RABO0000000000",
				Code:             "ID",
				AfBij:            "Bij",
				BedragEUR:        "22,00",
				Mutatiesoort:     "iDEAL",
				Mededelingen:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse id mauris at diam euismod maximus eget nec velit. Nulla eu scelerisque urna, ultricies varius nunc. Ut auctor velit id sodales sed.",
				SaldoNaMutatie:   "20,00",
			},
			AccountID: "test_accountid",
			PayeeName: "Shopping",
			Memo:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse id mauris at diam euismod maximus eget nec velit. Nulla eu scelerisque urna, ultricies varius nunc. Ut auctor velit id sodales",
			FlagColor: transaction.FlagColorGreen,
			Date:      time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			ImportID:  "YNAB:2200:2020-01-02:1",
			Amount:    22000,
		},
	}
	for _, test := range tests {
		expectedTrans := &transaction.PayloadTransaction{
			AccountID: test.AccountID,
			Date: api.Date{
				Time: test.Date,
			},
			Amount:    test.Amount,
			Cleared:   transaction.ClearingStatusCleared,
			Approved:  false,
			PayeeName: &test.PayeeName,
			Memo:      &test.Memo,
			FlagColor: &test.FlagColor,
			ImportID:  &test.ImportID,
		}

		trans, err := test.Line.ToYNAB(test.AccountID)
		t.Require().NoError(err)
		t.Require().Equal(expectedTrans, trans)
	}
}

func TestCSVTestSuite(t *testing.T) {
	suite.Run(t, new(CSVTestSuite))
}
