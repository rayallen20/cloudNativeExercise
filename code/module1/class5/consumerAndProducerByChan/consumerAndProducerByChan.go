package main

import (
	"fmt"
	"time"
)

// producer与consumer(worker)通信用
var messages = make(chan int, 10)

// 用于标识worker是否退出用
var done = make(chan bool)

func main() {
	go consumer()
	producer()
	defer close(messages)

	time.Sleep(5 * time.Second)
	close(done)
	time.Sleep(time.Second)
	fmt.Println("main process exit!")
}

func consumer() {
	ticker := time.NewTicker(time.Second)
	// 每秒消费1个
	for _ = range ticker.C {
		select {
		// close(done)会触发该分支
		case <-done:
			fmt.Println("consumer process interrupt...")
			return
		default:
			// Tip:此处原文写的是fmt.Printf("send message: %d\n", <-messages)
			// 个人推测可能写错了?
			fmt.Printf("received message:%d\n", <-messages)
		}
	}
}

func producer() {
	for i := 0; i < 10; i++ {
		messages <- i
	}
}
