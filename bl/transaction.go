package bl

import (
	"HDFC/poc/models"
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/jung-kurt/gofpdf"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdffont "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/font"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

type BL interface {
	WriteToPDF(transactions []models.Transaction)
}

type GoPDFWriter struct {
}

type PDFCPUWriter struct {
}

func GetTransactions() []models.Transaction {
	rand.Seed(time.Now().UnixNano())
	transactions := make([]models.Transaction, 1)

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
func (w GoPDFWriter) WriteToPDF(transactions []models.Transaction) {
	// Create a new PDF file with portrait orientation and millimeter units
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set up the table header
	header := []string{"ID", "Account ID", "Amount", "Type", "Date", "Description"}
	pdf.SetFont("Arial", "B", 16)
	for _, str := range header {
		pdf.CellFormat(35, 7, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Add the transactions to the table
	pdf.SetFont("Arial", "", 12)

	for _, t := range transactions {
		pdf.CellFormat(35, 7, strconv.FormatInt(t.ID, 10), "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 7, strconv.FormatInt(t.AccountNumber, 10), "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 7, strconv.FormatFloat(t.Amount, 'f', 2, 64), "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 7, string(t.Type), "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 7, t.Date.Format("2006-01-02 15:04:05"), "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 7, t.Description, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	// Save the PDF file to disk
	err := pdf.OutputFileAndClose("transactions.pdf")
	if err != nil {
		panic(err)
	}
}

// using pdfCPU
func (w PDFCPUWriter) WriteToPDF(transactions []models.Transaction) {
	// b, _ := json.Marshal(transactions)

	createPDF(transactions)
	// zero := 0
	// var offset int64 = 0
	// err := api.CreatePDFFile(&model.XRefTable{
	// 	Table: map[int]*model.XRefTableEntry{
	// 		1: &model.XRefTableEntry{
	// 			Free:       true,
	// 			Generation: &zero,
	// 			Offset:     &offset,
	// 		},
	// 	},
	// }, "t.pdf", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// // Save the PDF file to disk
	// err = api.OptimizeFile("t.pdf", "t.pdf", nil)
	// if err != nil {
	// 	panic(err)
	// }

}

// Helper function to create a PDF file with a table of transactions
func createPDF(transactions []models.Transaction) {
	// Create a new PDF context
	// _, err := model.NewContext(nil, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// Define the PDF page size and margins
	// pageSize := types.PaperSize["A4"]

	// mediaBox := types.RectForDim(pageSize.Width, pageSize.Height)

	// margins := primitives.Margin{Left: 10, Top: 10, Right: 10, Bottom: 10}

	// // Define the PDF font and font size
	// fontName := "Helvetica"
	// fontSize := 12

	// Create a new PDF page
	// page := model.Page{MediaBox: mediaBox}

	// Create a new table with header and rows
	// header := []string{"ID", "Account ID", "Amount", "Type", "Date", "Description"}
	msg := "TestCreateDemoPDF"
	mediaBox := types.RectForFormat("A4")
	p := model.Page{MediaBox: mediaBox, Fm: model.FontMap{}, Buf: new(bytes.Buffer)}
	CreateTestPageContent(p)
	xRefTable, err := CreateDemoXRef()
	if err != nil {
		// t.Fatalf("%s: %v\n", msg, err)
		panic(err)
	}
	rootDict, err := xRefTable.Catalog()
	if err != nil {
		panic(err)
	}
	if err = AddPageTreeWithSamplePage(xRefTable, rootDict, p); err != nil {
		panic(err)
	}
	createAndValidate(nil, xRefTable, "Test.pdf", msg)

	// rows := make([][]string, len(transactions))
	// for i, t := range transactions {
	// 	rows[i] = []string{
	// 		fmt.Sprintf("%d", t.ID),
	// 		fmt.Sprintf("%d", t.AccountNumber),
	// 		fmt.Sprintf("%.2f", t.Amount),
	// 		string(t.Type),
	// 		t.Date.String(),
	// 		t.Description,
	// 	}
	// }

	// xRefTable := &model.XRefTable{
	// 	Table:      map[int]*model.XRefTableEntry{},
	// 	Names:      map[string]*model.Node{},
	// 	PageAnnots: map[int]model.PgAnnots{},
	// 	Stats:      model.NewPDFStats(),
	// 	URIs:       map[int]map[string]string{},
	// 	UsedGIDs:   map[string]map[uint16]bool{},
	// }

	// xRefTable.Table[0] = model.NewFreeHeadXRefTableEntry()

	// one := 1
	// xRefTable.Size = &one

	// v := model.V17
	// xRefTable.HeaderVersion = &v

	// xRefTable.PageCount = 0

	// // Optional infoDict.
	// xRefTable.Info = nil

	// // Additional streams not implemented.
	// xRefTable.AdditionalStreams = nil

	// rootDict := types.NewDict()
	// rootDict.InsertName("Type", "Catalog")

	// ir, err := xRefTable.IndRefForNewObject(rootDict)
	// if err != nil {
	// 	panic(err)
	// }

	// xRefTable.Root = ir

	// err = api.CreatePDFFile(xRefTable, "t.pdf", nil)
	// if err != nil {
	// 	panic(err)
	// }
	// table := model.Table{
	// 	Header:  header,
	// 	Rows:    rows,
	// 	Style:   pdfcpu.NewTableStyle(fontSize, fontName, pdfcpu.Black),
	// 	Align:   pdfcpu.AlignCenter,
	// 	ColSize: pdfcpu.NewUSet(40, 40, 40, 40, 60, 100),
	// }

	// Add the table to the PDF page
	// if err := pdfcpu.CreateTable(&page, margins, table); err != nil {
	// 	return nil, err
	// }

	// // Add the PDF page to the PDF context
	// if err := ctx.AddPage(page); err != nil {
	// 	return nil, err
	// }

	// // Optimize the PDF document
	// if err := api.Optimize(ctx, nil); err != nil {
	// 	return nil, err
	// }

	// // Return the PDF document as a byte slice
	// pdf, err := api.ReadSeeker(ctx.Write.ReadSeeker())
	// if err != nil {
	// 	return nil, err
	// }
	// return pdf, nil
}

// CreateTestPageContent draws a test grid.
func CreateTestPageContent(p model.Page) {
	b := p.Buf
	mb := p.MediaBox

	b.WriteString("[3]0 d 0 w ")

	// X
	fmt.Fprintf(b, "0 0 m %f %f l s %f 0 m 0 %f l s ",
		mb.Width(), mb.Height(), mb.Width(), mb.Height())

	// Horizontal guides
	c := 6
	if mb.Landscape() {
		c = 4
	}
	j := mb.Height() / float64(c)
	for i := 1; i < c; i++ {
		k := mb.Height() - float64(i)*j
		s := fmt.Sprintf("0 %f m %f %f l s ", k, mb.Width(), k)
		b.WriteString(s)
	}

	// Vertical guides
	c = 4
	if mb.Landscape() {
		c = 6
	}
	j = mb.Width() / float64(c)
	for i := 1; i < c; i++ {
		k := float64(i) * j
		s := fmt.Sprintf("%f 0 m %f %f l s ", k, k, mb.Height())
		b.WriteString(s)
	}
}

func addContents(xRefTable *model.XRefTable, pageDict types.Dict, p model.Page) error {
	CreateTestPageContent(p)
	sd, _ := xRefTable.NewStreamDictForBuf(p.Buf.Bytes())

	if err := sd.Encode(); err != nil {
		return err
	}

	ir, err := xRefTable.IndRefForNewObject(*sd)
	if err != nil {
		return err
	}

	pageDict.Insert("Contents", *ir)

	return nil
}

func CreateXRefTableWithRootDict() (*model.XRefTable, error) {
	xRefTable := &model.XRefTable{
		Table:      map[int]*model.XRefTableEntry{},
		Names:      map[string]*model.Node{},
		PageAnnots: map[int]model.PgAnnots{},
		Stats:      model.NewPDFStats(),
		URIs:       map[int]map[string]string{},
		UsedGIDs:   map[string]map[uint16]bool{},
	}

	xRefTable.Table[0] = model.NewFreeHeadXRefTableEntry()

	one := 1
	xRefTable.Size = &one

	v := model.V17
	xRefTable.HeaderVersion = &v

	xRefTable.PageCount = 0

	// Optional infoDict.
	xRefTable.Info = nil

	// Additional streams not implemented.
	xRefTable.AdditionalStreams = nil

	rootDict := types.NewDict()
	rootDict.InsertName("Type", "Catalog")

	ir, err := xRefTable.IndRefForNewObject(rootDict)
	if err != nil {
		return nil, err
	}

	xRefTable.Root = ir

	return xRefTable, nil
}

// CreateDemoXRef creates a minimal single page PDF file for demo purposes.
func CreateDemoXRef() (*model.XRefTable, error) {
	xRefTable, err := CreateXRefTableWithRootDict()
	if err != nil {
		return nil, err
	}

	return xRefTable, nil
}

func AddPageTreeWithSamplePage(xRefTable *model.XRefTable, rootDict types.Dict, p model.Page) error {

	// mediabox = physical page dimensions
	mba := p.MediaBox.Array()

	pagesDict := types.Dict(
		map[string]types.Object{
			"Type":     types.Name("Pages"),
			"Count":    types.Integer(1),
			"MediaBox": mba,
		},
	)

	parentPageIndRef, err := xRefTable.IndRefForNewObject(pagesDict)
	if err != nil {
		return err
	}

	pageIndRef, err := createDemoPage(xRefTable, *parentPageIndRef, p)
	if err != nil {
		return err
	}

	pagesDict.Insert("Kids", types.Array{*pageIndRef})
	rootDict.Insert("Pages", *parentPageIndRef)

	return nil
}

func createDemoPage(xRefTable *model.XRefTable, parentPageIndRef types.IndirectRef, p model.Page) (*types.IndirectRef, error) {

	pageDict := types.Dict(
		map[string]types.Object{
			"Type":   types.Name("Page"),
			"Parent": parentPageIndRef,
		},
	)

	fontRes, err := pdffont.FontResources(xRefTable, p.Fm)
	if err != nil {
		return nil, err
	}

	if len(fontRes) > 0 {
		resDict := types.Dict(
			map[string]types.Object{
				"Font": fontRes,
			},
		)
		pageDict.Insert("Resources", resDict)
	}

	ir, err := createDemoContentStreamDict(xRefTable, pageDict, p.Buf.Bytes())
	if err != nil {
		return nil, err
	}
	pageDict.Insert("Contents", *ir)

	return xRefTable.IndRefForNewObject(pageDict)
}

func createDemoContentStreamDict(xRefTable *model.XRefTable, pageDict types.Dict, b []byte) (*types.IndirectRef, error) {
	sd, _ := xRefTable.NewStreamDictForBuf(b)
	if err := sd.Encode(); err != nil {
		return nil, err
	}
	return xRefTable.IndRefForNewObject(*sd)
}

func createAndValidate(t *testing.T, xRefTable *model.XRefTable, outFile, msg string) {
	// t.Helper()
	outDir := "."
	outFile = filepath.Join(outDir, outFile)
	if err := api.CreatePDFFile(xRefTable, outFile, nil); err != nil {
		// t.Fatalf("%s: %v\n", msg, err)
		panic(err)
	}
	if err := api.ValidateFile(outFile, nil); err != nil {
		// t.Fatalf("%s: %v\n", msg, err)
		panic(err)
	}
}
