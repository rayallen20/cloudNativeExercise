package main

import "fmt"

const MAX = 3

var numbers []int = []int{10, 100, 200}

func main() {
	correct()
	mistake()
}

func correct() {
	// 指针数组
	var ptrs [MAX]*int

	// 将numbers数组中的每一个元素的地址
	// 赋值给指针数组中的每一个元素
	for i := 0; i < MAX; i++ {
		ptrs[i] = &numbers[i]
	}

	for key, value := range ptrs {
		// 每个元素的内存地址是不同的
		fmt.Printf("pointerArray: index:%d value:%d, address of value:%d\n", key, *value, value)
	}
}

func mistake() {
	// 指针数组
	var ptrs [MAX]*int

	// 使用for range 将numbers数组中的每一个元素的地址
	// 赋值给指针数组中的每一个元素
	for key, value := range numbers {
		ptrs[key] = &value
	}

	for key, value := range ptrs {
		// 每个元素地址相同
		fmt.Printf("pointerArray: index:%d value:%d, address of value:%d\n", key, *value, value)
	}
}
