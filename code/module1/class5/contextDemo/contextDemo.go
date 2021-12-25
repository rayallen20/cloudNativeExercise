package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	baseCtx := context.Background()
	// 向上下文中写入一个键值对
	ctx := context.WithValue(baseCtx, "key", "value")
	go func(c context.Context) {
		// 被主线程创建的goroutine 通过上下文读取传递的变量
		fmt.Println(c.Value("key"))
	}(ctx)

	// 超时控制
	timeoutCtx, cancel := context.WithTimeout(baseCtx, time.Second)
	defer cancel()
	go func(ctx context.Context) {
		// 定时器 每秒向ticker.C中放入一个Time对象
		ticker := time.NewTicker(time.Second)
		// 每秒循环1次
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				// 超时处理
				fmt.Println("child process interrupt...")
				return
			default:
				// 正常逻辑
				fmt.Println("enter default")
			}
		}
	}(timeoutCtx)

	select {
	// 等待子线程完成
	case <-timeoutCtx.Done():
		time.Sleep(time.Second)
		fmt.Println("main process exit!")
	}
}
