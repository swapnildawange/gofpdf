package main

import (
	"HDFC/poc/bl"
	"fmt"
	"html"
	"log"
	"net/http"
)

const (
	PORT = 8080
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/hi")
		ts := bl.GetTransactions()

		gopdfWriter := &bl.GoPDFWriter{}
		pdfBase64 := gopdfWriter.WriteToPDF(ts)

		w.Write([]byte(pdfBase64))
	})

	errChan := make(chan error, 1)
	go func() {
		errChan <- http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	}()
	log.Printf("running REST server on PORT : %d \n", PORT)
	<-errChan
}
