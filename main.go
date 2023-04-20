package main

import (
	"HDFC/poc/bl"
	"fmt"
	"log"
	"net/http"
)

const (
	PORT = 8080
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ts := bl.GetTransactions(100)
		gopdfWriter := &bl.GoPDFWriter{}
		pdfBase64 := gopdfWriter.WriteToPDF(ts)
		w.Write([]byte(pdfBase64))
	})

	errChan := make(chan error, 1)
	go func() {
		errChan <- http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	}()
	log.Printf("running REST server on PORT : %d \n", PORT)
	log.Println(<-errChan)

}
