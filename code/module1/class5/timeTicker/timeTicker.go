package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	timer := time.NewTicker(time.Second)
	ch := make(chan int)
	select {
	case <-ch:
		fmt.Println("received from ch")
	case <-timer.C:
		fmt.Println("timeout waiting from channel ch")
	}
	context.Background()
}
