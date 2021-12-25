package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			// generate a random int n and n will be between 0 and 10
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(10)
			fmt.Println("putting:", n)
			ch <- n
		}
		close(ch)
	}()

	fmt.Println("hello from main")
	for v := range ch {
		fmt.Println("receiving:", v)
	}
}
