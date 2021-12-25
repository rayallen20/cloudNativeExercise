package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	defer fmt.Printf("1\n")
	defer fmt.Printf("2\n")
	defer fmt.Printf("3\n")
	loopFunc()
	time.Sleep(1 * time.Second)
}

func loopFunc() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		go func(i int) {
			lock.Lock()
			defer lock.Unlock()
			fmt.Println("loopFunc:", i)
		}(i)
	}
}
