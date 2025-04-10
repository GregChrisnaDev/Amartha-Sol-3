package usecase

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/pdfgenerator"
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

type CountROIProfitReq struct {
	Rate            float64
	PrincipalAmount float64
	LoanDuration    float64
	LendAmount      float64
}

func countROIProfit(params CountROIProfitReq) (float64, float64) {
	ratePerWeek := float64(params.Rate) / 52
	totalInterest := params.PrincipalAmount * ratePerWeek * float64(params.LoanDuration) / 100

	profit := params.LendAmount / params.PrincipalAmount * totalInterest
	roi := profit / params.LendAmount * 100

	return roi, profit
}

type GenerateAgreementLetterReq struct {
	NameLender    string
	NameLoaner    string
	AddressLender string
	AddressLoaner string
	SignLender    string
	SignLoaner    string
	Filename      string

	CountROIProfitReq
}

func generateAgreementLetter(client pdfgenerator.Client, params GenerateAgreementLetterReq) error {
	_, profit := countROIProfit(CountROIProfitReq{
		Rate:            params.Rate,
		PrincipalAmount: params.PrincipalAmount,
		LoanDuration:    params.LoanDuration,
		LendAmount:      params.LendAmount,
	})

	return client.GenerateAgreementLetter(pdfgenerator.AgreementLetterPDF{
		NameLender:    params.NameLender,
		NameLoaner:    params.NameLoaner,
		AddressLender: params.AddressLender,
		AddressLoaner: params.AddressLoaner,
		SignLender:    params.SignLender,
		SignLoaner:    params.SignLoaner,
		Filename:      params.Filename,
		LendAmount:    convertToCurrency(params.LendAmount),
		LoanPayAmount: convertToCurrency(params.LendAmount + profit),
		LoanDuration:  uint32(params.LoanDuration),
	})
}
