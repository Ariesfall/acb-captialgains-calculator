package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

func main() {
	http.HandleFunc("/upload", csvHandler)
	log.Println("starting service... listen port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func csvHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receive incoming request %s - %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)

	// Parse the CSV file from the request body.
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
