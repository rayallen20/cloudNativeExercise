package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("done channel is triggered, exit child goroutine")
				return
			default:
				fmt.Println("handle biz")
			}
		}
	}()

	time.Sleep(time.Millisecond)
	close(done)
}
