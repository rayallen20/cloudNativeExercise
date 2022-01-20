package main

func main() {
	// 子线程消耗1个CPU
	go func() {
		for {
		}
	}()
	// 主线程消耗1个CPU
	for {
	}
}
