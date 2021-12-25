package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		fmt.Println("hello from goroutine")
		ch <- 0
	}()
	<-ch
}
