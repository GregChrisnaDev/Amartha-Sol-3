package pdfgenerator

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"codeberg.org/go-pdf/fpdf"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/storage"
)

type TemplateData struct {
	Title       string
	ImageBase64 template.HTML
}

func (c *client) GenerateAgreementLetter(data AgreementLetterPDF) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(0, 10, "PERJANJIAN PINJAMAN")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	// Intro
	intro := fmt.Sprintf(`Perjanjian Pinjaman ini ("Perjanjian") dibuat dan ditandatangani, oleh dan antara:

		%s, berkedudukan di %s.

		%s, berkedudukan di %s.

		Kreditur dan Debitur selanjutnya secara bersama-sama disebut "Para Pihak".`, data.NameLender, data.AddressLender, data.NameLoaner, data.AddressLoaner)
	pdf.MultiCell(0, 7, intro, "", "", false)
	pdf.Ln(5)

	// Pasal sections
	articles := []string{
		"Pasal 1 - PINJAMAN",
		"Pasal 2 - JAMINAN",
		"Pasal 3 - BUNGA",
		"Pasal 4 - JANGKA WAKTU PERJANJIAN",
		"Pasal 5 - PELUNASAN PINJAMAN",
		"Pasal 6 - JAMINAN DEBITUR",
		"Pasal 7 - PAJAK DAN BIAYA",
		"Pasal 8 - PERNYATAAN DAN JAMINAN",
		"Pasal 9 - KEADAAN CIDERA JANJI",
		"Pasal 10 - PEMBERITAHUAN",
		"Pasal 11 - PENGALIHAN",
		"Pasal 12 - PILIHAN HUKUM DAN DOMISILI",
		"Pasal 13 - PENYELESAIAN SENGKETA",
		"Pasal 14 - KERAHASIAAN",
		"Pasal 15 - KETERPISAHAN",
		"Pasal 16 - KESELURUHAN PERJANJIAN",
		"Pasal 17 - HAL-HAL LAIN",
	}

	for _, article := range articles {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 7, article)
		pdf.Ln(7)
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(0, 6, "...", "", "", false)
		pdf.Ln(3)
	}

	// Add some space before the signature row
	pdf.Ln(10)

	// Remember current Y position
	y := pdf.GetY()

	// Width for each signature block (half of page width minus some margin)
	halfWidth := 90.0

	// Kreditur (left side)
	pdf.SetY(y)
	pdf.SetX(10) // margin from left
	pdf.CellFormat(halfWidth, 10, "Kreditur", "", 0, "C", false, 0, "")
	pdf.Ln(20)
	pdf.Ln(7)
	pdf.ImageOptions(c.storagePath+data.SignLender, 30, y+10, 50, 35, false, fpdf.ImageOptions{
		ImageType: "JPG",
	}, 0, "")
	pdf.Ln(20)
	pdf.SetX(10)
	pdf.CellFormat(halfWidth, 10, data.NameLender, "", 0, "C", false, 0, "")

	// Debitur (right side)
	pdf.SetY(y)
	pdf.SetX(120) // half of A4 width + margin
	pdf.CellFormat(halfWidth, 10, "Debitur", "", 0, "C", false, 0, "")
	pdf.Ln(20)

	pdf.Ln(7)
	if data.SignLoaner != "" {
		pdf.ImageOptions(c.storagePath+data.SignLoaner, 140, y+10, 50, 35, false, fpdf.ImageOptions{
			ImageType: "JPG",
		}, 0, "")

	}
	pdf.Ln(20)
	pdf.SetX(120)
	pdf.CellFormat(halfWidth, 10, data.NameLoaner, "", 0, "C", false, 0, "")

	if err := os.MkdirAll(c.storagePath+storage.AGREEMENT_LETTER_DIR, os.ModePerm); err != nil {
		log.Println("[GeneratePDF] error while mkdir", err.Error())
		return err
	}

	// Save file
	err := pdf.OutputFileAndClose(c.storagePath + storage.AGREEMENT_LETTER_DIR + data.Filename)
	if err != nil {
		log.Println("[GeneratePDF] error write", err.Error())
		return err
	}

	return nil
}
