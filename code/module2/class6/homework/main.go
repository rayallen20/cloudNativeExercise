package main

import (
	"context"
	"flag"
	"homework2/controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var listenAddr string

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// 通知主线程优雅退出工作已完成的chan
	done := make(chan bool, 1)

	// 优雅退出goroutine从OS收集signal的chan
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := newWebServer(logger)
	go gracefullyShutdown(server, logger, quit, done)

	logger.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func gracefullyShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan bool) {
	<-quit
	logger.Println("Serveer is shutting down")

	// TODO: 其实这个30应该也用flag解析 太晕乎了不干了 就这样凑合凑合吧
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server:%v\n", err)
	}

	close(done)
}

func newWebServer(logger *log.Logger) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/getHeader", controller.GetHeader)
	router.HandleFunc("/healthz", controller.Healthz)

	return &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
