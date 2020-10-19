package csv

import "go.bmvs.io/ynab/api/transaction"

type TransactionLine interface {
	ToYnab() *transaction.PayloadTransaction
}
