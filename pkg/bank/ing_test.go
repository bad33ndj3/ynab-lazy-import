package bank

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
)

type INGTestSuite struct {
	suite.Suite
}

func (t *INGTestSuite) TestToYNAB() {
	tests := []struct {
		Line      ing
		AccountID string
		PayeeName string
		Memo      string
		FlagColor transaction.FlagColor
		Date      time.Time
		ImportID  string
		Amount    int64
	}{
		{
			Line: ing{
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
			Line: ing{
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
			Line: ing{
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
			Line: ing{
				Datum:            20200102,
				NaamOmschrijving: "Shopping",
				Rekening:         "NL13INGB0000000000",
				Tegenrekening:    "NL14RABO0000000000",
				Code:             "ID",
				AfBij:            "Bij",
				BedragEUR:        "22,00",
				Mutatiesoort:     "iDEAL",
				Mededelingen: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse id mauris at " +
					"diam euismod maximus eget nec velit. Nulla eu scelerisque urna, ultricies varius nunc. Ut auctor " +
					"velit id sodales sed.",
				SaldoNaMutatie: "20,00",
			},
			AccountID: "test_accountid",
			PayeeName: "Shopping",
			Memo: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse id mauris at diam euismod" +
				" maximus eget nec velit. Nulla eu scelerisque urna, ultricies varius nunc. Ut auctor velit id sodales",
			FlagColor: transaction.FlagColorGreen,
			Date:      time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			ImportID:  "YNAB:2200:2020-01-02:1",
			Amount:    22000,
		},
	}
	for index := range tests {
		expectedTrans := &transaction.PayloadTransaction{
			AccountID: tests[index].AccountID,
			Date: api.Date{
				Time: tests[index].Date,
			},
			Amount:    tests[index].Amount,
			Cleared:   transaction.ClearingStatusCleared,
			Approved:  false,
			PayeeName: &tests[index].PayeeName,
			Memo:      &tests[index].Memo,
			FlagColor: &tests[index].FlagColor,
			ImportID:  &tests[index].ImportID,
		}

		trans, err := tests[index].Line.toYNAB(tests[index].AccountID)
		t.Require().NoError(err)
		t.Require().Equal(expectedTrans, trans)
	}
}

func TestINGTestSuite(t *testing.T) {
	suite.Run(t, new(INGTestSuite))
}
