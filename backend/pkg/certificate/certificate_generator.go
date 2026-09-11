package certificate

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
)

type CertificateData struct {
	ParticipantName string
	ParticipantNIM  string
	University      string
	PeriodTitle     string
	CertificateExp  string
	ListeningScore  int
	StructureScore  int
	ReadingScore    int
	FinalScore      int
	TemplateURL     string
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

// Generate creates a PDF certificate based on the template (PDF or Image)
// and overlays student data, scores, period, and expiration date.
func (g *Generator) Generate(data *CertificateData) ([]byte, error) {
	// A4 Landscape: 297mm x 210mm
	const pageWidth = 297.0
	const pageHeight = 210.0

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	var templateApplied bool

	if data.TemplateURL != "" {
		templateBytes, contentType, err := downloadTemplate(data.TemplateURL)
		if err == nil && len(templateBytes) > 0 {
			if isPDF(contentType, data.TemplateURL) {
				// Apply PDF template via gofpdi
				err = applyPDFTemplate(pdf, templateBytes, pageWidth, pageHeight)
				if err == nil {
					templateApplied = true
				}
			} else {
				// Apply Image template (JPEG/PNG)
				err = applyImageTemplate(pdf, templateBytes, contentType, pageWidth, pageHeight)
				if err == nil {
					templateApplied = true
				}
			}
		}
	}

	// If no template or failed to download template, draw default background
	if !templateApplied {
		drawDefaultBackground(pdf, pageWidth, pageHeight)
	}

	// Overlay dynamic certificate text onto the template
	overlayCertificateDetails(pdf, data, pageWidth)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate certificate pdf: %w", err)
	}

	return buf.Bytes(), nil
}

func downloadTemplate(url string) ([]byte, string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download template, status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	contentType := resp.Header.Get("Content-Type")
	return data, contentType, nil
}

func isPDF(contentType, url string) bool {
	if strings.Contains(strings.ToLower(contentType), "pdf") {
		return true
	}
	ext := strings.ToLower(filepath.Ext(url))
	return ext == ".pdf"
}

func applyPDFTemplate(pdf *gofpdf.Fpdf, templateBytes []byte, width, height float64) error {
	tmpFile, err := os.CreateTemp("", "cert_template_*.pdf")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(templateBytes); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	tpl := gofpdi.ImportPage(pdf, tmpFile.Name(), 1, "/MediaBox")
	gofpdi.UseImportedTemplate(pdf, tpl, 0, 0, width, height)
	return nil
}

func applyImageTemplate(pdf *gofpdf.Fpdf, imgBytes []byte, contentType string, width, height float64) error {
	imgType := "PNG"
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		imgType = "JPG"
	}

	imgName := "cert_bg_" + fmt.Sprintf("%d", time.Now().UnixNano())
	opt := gofpdf.ImageOptions{
		ImageType: imgType,
		ReadDpi:   true,
	}

	pdf.RegisterImageOptionsReader(imgName, opt, bytes.NewReader(imgBytes))
	pdf.ImageOptions(imgName, 0, 0, width, height, false, opt, 0, "")
	return nil
}

func drawDefaultBackground(pdf *gofpdf.Fpdf, width, height float64) {
	// Outer border
	pdf.SetDrawColor(41, 65, 114) // Deep navy
	pdf.SetLineWidth(2.5)
	pdf.Rect(10, 10, width-20, height-20, "D")

	// Inner border
	pdf.SetLineWidth(0.8)
	pdf.SetDrawColor(180, 150, 90) // Gold
	pdf.Rect(14, 14, width-28, height-28, "D")

	// Header
	pdf.SetTextColor(41, 65, 114)
	pdf.SetFont("Arial", "B", 26)
	pdf.SetXY(0, 32)
	pdf.CellFormat(width, 12, "CERTIFICATE OF ACHIEVEMENT", "", 0, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(100, 100, 100)
	pdf.SetXY(0, 48)
	pdf.CellFormat(width, 8, "This is to certify that", "", 0, "C", false, 0, "")
}

func overlayCertificateDetails(pdf *gofpdf.Fpdf, data *CertificateData, width float64) {
	// 1. Participant Name
	pdf.SetTextColor(30, 41, 59) // Dark Slate
	pdf.SetFont("Arial", "B", 24)
	pdf.SetXY(0, 72)
	pdf.CellFormat(width, 12, strings.ToUpper(data.ParticipantName), "", 0, "C", false, 0, "")

	// 2. Participant NIM & University
	pdf.SetTextColor(71, 85, 105)
	pdf.SetFont("Arial", "", 12)
	infoText := ""
	if data.ParticipantNIM != "" && data.University != "" {
		infoText = fmt.Sprintf("NIM: %s  |  %s", data.ParticipantNIM, data.University)
	} else if data.ParticipantNIM != "" {
		infoText = fmt.Sprintf("NIM: %s", data.ParticipantNIM)
	} else if data.University != "" {
		infoText = data.University
	}
	if infoText != "" {
		pdf.SetXY(0, 86)
		pdf.CellFormat(width, 7, infoText, "", 0, "C", false, 0, "")
	}

	// 3. Period Title
	if data.PeriodTitle != "" {
		pdf.SetFont("Arial", "I", 12)
		pdf.SetTextColor(100, 116, 139)
		pdf.SetXY(0, 96)
		pdf.CellFormat(width, 7, fmt.Sprintf("has completed the %s", data.PeriodTitle), "", 0, "C", false, 0, "")
	}

	// 4. Section Scores Table
	// Table layout: 3 columns centered
	scoreTableWidth := 180.0
	startX := (width - scoreTableWidth) / 2.0
	colWidth := scoreTableWidth / 3.0
	tableY := 115.0

	// Header background & border
	pdf.SetFillColor(241, 245, 249)
	pdf.SetDrawColor(203, 213, 225)
	pdf.SetLineWidth(0.3)
	pdf.SetTextColor(51, 65, 85)
	pdf.SetFont("Arial", "B", 10)

	pdf.SetXY(startX, tableY)
	pdf.CellFormat(colWidth, 9, "Listening Comprehension", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidth, 9, "Structure & Written Expr.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidth, 9, "Reading Comprehension", "1", 0, "C", true, 0, "")

	// Score Values
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetXY(startX, tableY+9)
	pdf.CellFormat(colWidth, 12, fmt.Sprintf("%d", data.ListeningScore), "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidth, 12, fmt.Sprintf("%d", data.StructureScore), "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidth, 12, fmt.Sprintf("%d", data.ReadingScore), "1", 0, "C", false, 0, "")

	// 5. Total Final Score Box
	totalBoxY := tableY + 26.0
	pdf.SetFillColor(238, 242, 255)
	pdf.SetDrawColor(99, 102, 241)
	pdf.SetLineWidth(0.5)
	totalWidth := 100.0
	totalX := (width - totalWidth) / 2.0

	pdf.Rect(totalX, totalBoxY, totalWidth, 16, "FD")
	pdf.SetFont("Arial", "B", 13)
	pdf.SetTextColor(67, 56, 202)
	pdf.SetXY(totalX, totalBoxY)
	pdf.CellFormat(totalWidth, 16, fmt.Sprintf("TOTAL TOEFL SCORE: %d", data.FinalScore), "", 0, "C", false, 0, "")

	// 6. Expiry Date & Issue Notice
	if data.CertificateExp != "" {
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(148, 163, 184)
		pdf.SetXY(0, totalBoxY+22)
		pdf.CellFormat(width, 6, fmt.Sprintf("Valid Until: %s", data.CertificateExp), "", 0, "C", false, 0, "")
	}
}

