package pkg

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
)

func CsvHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receive incoming request %s - %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
	log.Println(r.Body)
	reader := csv.NewReader(r.Body)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(records)

	result := Calculate(records)

	// Return the response in the desired format.
	var jsonBytes []byte
	if reflect.TypeOf(result) != reflect.TypeOf(jsonBytes) {
		jsonBytes, err = json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		jsonBytes = result
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func CsvReader(filename string) {
	log.Println("filename: ", filename)
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}

	reader := csv.NewReader(bytes.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
		return
	}

	result := Calculate(records)
	log.Println(string(result))
}
