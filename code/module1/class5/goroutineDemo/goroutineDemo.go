package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go fmt.Println(i)
	}
	time.Sleep(time.Second)
}
