package bl

import (
	"testing"
)

func Test_WriteToPDF(t *testing.T) {
	w := GoPDFWriter{}
	transactions := GetTransactions(100)
	w.WriteToPDF(transactions)

}
