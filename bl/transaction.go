package bl

import (
	"HDFC/poc/models"
	"HDFC/poc/utils"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type BL interface {
	WriteToPDF(transactions []models.Transaction) string
}

type GoPDFWriter struct {
}

type PDFCPUWriter struct {
}

const (
	MARGIN = 10
)

func GetTransactions() []models.Transaction {
	rand.Seed(time.Now().UnixNano())
	transactions := make([]models.Transaction, 10000)

	for i := range transactions {
		transactions[i] = models.Transaction{
			ID:            rand.Int63(),
			AccountNumber: rand.Int63(),
			Amount:        float64(rand.Intn(1000)+1) / 10.0,
			Type:          getRandomType(),
			Date:          time.Now().Add(time.Duration(rand.Intn(86400)) * time.Second * -1),
			Description:   fmt.Sprintf("Transaction %d", i+1),
		}
	}

	return transactions
}

func getRandomType() models.TransactionType {
	types := []models.TransactionType{"Deposit", "Withdrawal", "Transfer"}
	return types[rand.Intn(len(types))]
}

// using gopdf
func (w GoPDFWriter) WriteToPDF(transactions []models.Transaction) string {
	// Create a new PDF file with portrait orientation and millimeter units
	pdf := gofpdf.New("P", "mm", "A4", "")
	// Set up the table title
	pdf.SetTitle("account-summary", true)

	// set footer
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})

	pdf.AliasNbPages("")

	columns := []string{"ID", "Account ID", "Amount", "Type", "Date", "Description"}

	pdf.AddPage()

	width, _ := pdf.GetPageSize()
	pdf.SetMargins(MARGIN, MARGIN, MARGIN)
	// remove the margins from total width
	// tw := (width - MARGIN*2) / float64(len(columns))
	ratios := []float64{0.15, 0.15, 0.15, 0.15, 0.2, 0.2}
	minWidth := 10.0

	columnWidths := utils.CalculateColumnWidths(ratios, width-2*MARGIN, minWidth)

	for j, str := range columns {
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(columnWidths[j], 10, str, "1", 0, "L", false, 0, "")
	}

	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)

	addTransactions(pdf, transactions, columnWidths)

	// Create a bytes.Buffer to use as an intermediary
	buffer := bytes.Buffer{}
	// Create an io.Writer that writes to the buffer
	writer := io.Writer(&buffer)
	// Write pdf data to the writer
	err := pdf.Output(writer)
	if err != nil {
		panic(err)
	}
	// Read the data from the buffer
	data := buffer.Bytes()
	// encode the data to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(data)
	return pdfBase64
}

func addTransactions(pdf *gofpdf.Fpdf, transactions []models.Transaction, cols []float64) {
	marginCell := 2.0
	_, pageh := pdf.GetPageSize()
	_, _, _, mbottom := pdf.GetMargins()
	rows := [][]string{}

	// create table rows
	for _, t := range transactions {
		rows = append(rows, []string{strconv.FormatInt(t.ID, 10), strconv.FormatInt(t.AccountNumber, 10),
			strconv.FormatFloat(t.Amount, 'f', 2, 64), string(t.Type), t.Date.String(), t.Description})
	}

	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx

		height := 0.
		_, lineHt := pdf.GetFontSize()

		for i, txt := range row {
			lines := pdf.SplitLines([]byte(txt), cols[i])
			h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
			if h > height {
				height = h
			}
		}
		// add a new page if the height of the row doesn't fit on the page
		if pdf.GetY()+height > pageh-mbottom {
			pdf.AddPage()
			y = pdf.GetY()
		}
		for i, txt := range row {
			width := cols[i]
			pdf.Rect(x, y, width, height, "")
			pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
			x += width
			pdf.SetXY(x, y)
		}
		pdf.SetXY(curx, y+height)
	}
}
