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
		log.Println("/")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(`[
  {
    "stateName": "Maharashtra",
    "cities": [
      {
        "cityName": "Ahmadpur",
        "branches": [
          {
            "branchName": "Ahmadpur Branch",
            "codBrn": 3877
          }
        ]
      },
      {
        "cityName": "Mumbai",
        "branches": [
          {
            "branchName": "Mumbai Branch 1",
            "codBrn": 1234
          },
          {
            "branchName": "Mumbai Branch 2",
            "codBrn": 5678
          }
        ]
      }
    ]
  },
  {
    "stateName": "Karnataka",
    "cities": [
      {
        "cityName": "Bengaluru",
        "branches": [
          {
            "branchName": "Bengaluru Branch 1",
            "codBrn": 9012
          },
          {
            "branchName": "Bengaluru Branch 2",
            "codBrn": 3456
          }
        ]
      },
      {
        "cityName": "Mangalore",
        "branches": [
          {
            "branchName": "Mangalore Branch 1",
            "codBrn": 7890
          }
        ]
      }
    ]
  }
]
`))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/hi")
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
	<-errChan
}
