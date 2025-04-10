package pdfgenerator

import "os"

type client struct {
	storagePath string
}

type Client interface {
	GenerateAgreementLetter(data AgreementLetterPDF) error
}

func Init() Client {
	storage := os.Getenv("DEFAULT_STORAGE")
	if storage == "" {
		storage = "./etc/storage/"
	}

	return &client{
		storagePath: storage,
	}
}
