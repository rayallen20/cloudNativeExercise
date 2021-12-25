package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt...")
				return
			default:
				fmt.Println("handle biz")
			}
		}
	}(timeoutCtx)

	time.Sleep(time.Microsecond)
	fmt.Println("err:", ctx.Err())
}
