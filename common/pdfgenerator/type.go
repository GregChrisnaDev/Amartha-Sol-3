package pdfgenerator

type AgreementLetterPDF struct {
	NameLender    string
	NameLoaner    string
	AddressLender string
	AddressLoaner string
	SignLender    string
	SignLoaner    string
	Filename      string
	LendAmount    string
	LoanPayAmount string
	LoanDuration  uint32
}
