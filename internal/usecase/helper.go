package usecase

import (
	"github.com/leekchan/accounting"
)

func convertToCurrency(value float64) string {
	ac := accounting.Accounting{
		Symbol:    "Rp",
		Precision: 2,
		Thousand:  ".",
		Decimal:   ",",
	}

	return ac.FormatMoney(value)
}
