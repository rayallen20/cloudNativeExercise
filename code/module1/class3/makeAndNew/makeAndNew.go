package main

import "fmt"

func main() {
	// 该指针中的地址 指向的切片长度和容量均为0
	mySlice1 := new([]int)
	fmt.Printf("len:%d, cap:%d\n", len(*mySlice1), cap(*mySlice1))

	// 以下是比较奇怪的用法 但是结果正确
	*mySlice1 = append(*mySlice1, 1)
	fmt.Printf("%#v\n", mySlice1)
	(*mySlice1)[0] = 2
	fmt.Printf("%#v\n", mySlice1)

	mySlice2 := make([]int, 0)
	mySlice3 := make([]int, 10)
	mySlice4 := make([]int, 10, 20)

	fmt.Printf("mySlice1: %p\n", mySlice1)
	fmt.Printf("mySlice2: len:%d cap:%d\n", len(mySlice2), cap(mySlice2))
	fmt.Printf("mySlice3: len:%d cap:%d\n", len(mySlice3), cap(mySlice3))
	fmt.Printf("mySlice4: len:%d cap:%d\n", len(mySlice4), cap(mySlice4))
}
