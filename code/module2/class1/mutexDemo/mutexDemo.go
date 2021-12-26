package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Tip:死锁定义
	// 死锁:所有并发进程彼此等待的程序.在这种情况下,如果没有外界的干预,这个程序将永远无法恢复

	go rLock()
	// Tip:直接调用wLock()会死锁
	go wLock()
	// Tip:直接调用lock()会死锁
	go lock()
	time.Sleep(5 * time.Second)
}

func lock() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		// Tip: 第2次循环就会阻塞 因为程序要等待lock解锁才能加锁
		// 但是给lock解锁的代码在循环期间又不会被执行 所以会阻塞
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("lock:", i)
	}
}

func rLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 3; i++ {
		// Tip:不会阻塞 因为RWMutex允许并发读
		lock.RLock()
		defer lock.RUnlock()
		fmt.Println("rLock:", i)
	}
}

func wLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 3; i++ {
		// Tip:会阻塞 因为RWMutex不允许并发写
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("wLock:", i)
	}
}
