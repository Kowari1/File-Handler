package pdf

import (
	"fmt"
	"path/filepath"
	"strconv"

	model "github.com/Kowari1/File-Handler/internal/domain"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type Device struct {
	N     int
	MsgID string
	Text  string
	Class int
	Level string
	Area  string
}

type Generator interface {
	Generate(outputDir string, unitGUID uuid.UUID, devices []Device) error
}

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) Generate(
	outputDir string,
	unitGUID uuid.UUID,
	devices []*model.Device,
) error {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 8)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Report for UnitGUID: %s", unitGUID))
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 7)

	headers := []string{
		"N", "MQTT", "Invid", "MsgID", "Text",
		"Context", "Class", "Level", "Area",
		"Addr", "Block", "Type", "Bit", "InvertBit",
	}

	colWidth := 15.0

	for _, h := range headers {
		pdf.CellFormat(colWidth, 6, h, "1", 0, "", false, 0, "")
	}
	pdf.Ln(-1)

	for _, d := range devices {

		row := []string{
			strconv.Itoa(d.N),
			d.MQTT,
			d.Invid,
			d.MsgID,
			d.Text,
			d.Context,
			strconv.Itoa(d.Class),
			d.Level,
			d.Area,
			d.Addr,
			d.Block,
			d.Type,
			d.Bit,
			d.InvertBit,
		}

		for _, cell := range row {
			pdf.CellFormat(colWidth, 6, cell, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}

	filePath := filepath.Join(outputDir, unitGUID.String()+".pdf")

	return pdf.OutputFileAndClose(filePath)
}
