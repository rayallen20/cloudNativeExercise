package main

import (
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.spreadsheetdb.io/")
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
