package pdfgenerator

type client struct {
	storagePath string
}

type Client interface {
	GenerateAgreementLetter(data AgreementLetterPDF) error
}

func Init(template, storage string) Client {

	if storage == "" {
		storage = "./etc/storage/"
	}

	return &client{
		storagePath: storage,
	}
}
