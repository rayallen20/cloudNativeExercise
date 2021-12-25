package main

import (
	"fmt"
	"time"
)

var message = make(chan int, 10)

func main() {
	go consumer()
	producer()
}

func consumer() {
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		select {
		case v := <-message:
			fmt.Println(v)
		}
	}
}

func producer() {
	for i := 0; i < 10; i++ {
		message <- i
		time.Sleep(time.Second)
	}
}
