package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ariesfall/acb-captialgains-calculator/pkg"
)

func main() {
	if len(os.Args) > 1 {
		pkg.CsvReader(os.Args[1])
		return
	}
	http.HandleFunc("/upload", pkg.CsvHandler)
	log.Println("starting service... listen port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
