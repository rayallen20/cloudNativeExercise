package main

import (
	"flag"
	"log"
	"os"
)

var listenAddr string

func main() {
	flag.StringVar(&listenAddr, "listen-addr", "5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is ready to handle requests at ,", listenAddr)
}
