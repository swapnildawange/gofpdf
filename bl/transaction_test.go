package bl

import (
	"testing"
)

func Benchmark_WriteToPDF(b *testing.B) {
	w := GoPDFWriter{}
	for n := 0; n < b.N; n++ {
		transactions := GetTransactions(n)
		w.WriteToPDF(transactions)
	}
}
