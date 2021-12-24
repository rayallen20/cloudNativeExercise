package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}

	fullString := "hello world"
	fmt.Println(fullString)
	for i, c := range fullString {
		// c的数据类型为rune 需转换为string 打印时才是这个字符的字面量
		// 否则打印的是ASCII码
		fmt.Println(i, string(c))
	}
}
