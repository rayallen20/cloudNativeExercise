package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var listenAddr string

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is ready to handle requests at ,", listenAddr)

	router := http.NewServeMux()
	router.HandleFunc("/", Demo)

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	logger.Println("Server stopped")
}

func Demo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
