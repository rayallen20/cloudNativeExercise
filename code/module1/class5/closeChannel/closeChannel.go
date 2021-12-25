package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)
	defer close(ch)
	if v, notClosed := <-ch; notClosed {
		fmt.Println(v)
	}
}
