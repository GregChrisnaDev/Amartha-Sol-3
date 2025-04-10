package usecase

import (
	"encoding/base64"
	"log"
	"os"

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

func encodeImgToBase64(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
}
