package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
	"testing"
	"time"
)

type CSVTestSuite struct {
	suite.Suite

	testAssets string
}

func (s *CSVTestSuite) SetupSuite() {
	s.testAssets = "testassets"
}

func (s *CSVTestSuite) TestCSVToINGExport() {
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
					Mededelingen:     "Naam: Origin.com EA Omschrijving: 0123455678 00100000123344556 EAPlay IBAN: NL14RABO0000000000 Kenmerk: 01-01-2020 09:00 0020004331516784 Valutadatum: 01-01-2020",
					SaldoNaMutatie:   "20,00",
					Tag:              "",
				},
			},
			inAccount: "NL13INGB0000000000",
			inDir:     fmt.Sprintf("%s/%s", s.testAssets, "base"),
		},
		{
			inAccount: "accountNonExisting",
			inDir:     fmt.Sprintf("%s/%s", s.testAssets, "base"),
		},
		{
			inAccount: "accountNonExisting",
			inDir:     fmt.Sprintf("%s/%s", s.testAssets, "pathNonExisting"),
			err:       errFailedToGetPath,
		},
	}
	for _, test := range tests {
		s.Run("", func() {
			lines, err := getLines(test.inAccount, test.inDir)
			s.Require().Equal(test.err, err)
			s.Require().Equal(test.expectedOutput, lines)
		})
	}
}

func (s CSVTestSuite) TestToYNAB() {
	tests := []struct {
		Line      INGExport
		AccountID string
		PayeeName string
		Memo      string
		FlagColor transaction.FlagColor
		Date      time.Time
		ImportID  string
	}{
		{
			Line: INGExport{
				Datum:            20200102,
				NaamOmschrijving: "Origin.com EA",
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
			PayeeName: "Origin.com EA",
			Memo:      "example",
			FlagColor: transaction.FlagColorGreen,
			Date:      time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC),
			ImportID:  "YNAB:-399:2020-01-02:1",
		},
	}
	for _, test := range tests {
		expectedTrans := &transaction.PayloadTransaction{
			AccountID: test.AccountID,
			Date: api.Date{
				Time: test.Date,
			},
			Amount:    -399,
			Cleared:   transaction.ClearingStatusCleared,
			Approved:  false,
			PayeeName: &test.PayeeName,
			Memo:      &test.Memo,
			FlagColor: &test.FlagColor,
			ImportID:  &test.ImportID,
		}

		trans, err := test.Line.ToYNAB(test.AccountID)
		s.Require().NoError(err)
		s.Require().Equal(expectedTrans, trans)
	}
}

func TestCSVTestSuite(t *testing.T) {
	suite.Run(t, new(CSVTestSuite))
}
