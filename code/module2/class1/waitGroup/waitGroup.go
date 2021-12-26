package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	waitByWG()
}

func waitBySleep() {
	for i := 0; i < 100; i++ {
		go fmt.Println(i)
	}
	time.Sleep(time.Second)
}

func waitByChannel() {
	// Tip:如果通过管道传递的数据仅做标识作用 通常使用空结构体
	c := make(chan struct{}, 100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			c <- struct{}{}
		}(i)
	}

	for i := 0; i < 100; i++ {
		// 从c中读不出100个元素 就阻塞
		<-c
	}
}

func waitByWG() {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}
