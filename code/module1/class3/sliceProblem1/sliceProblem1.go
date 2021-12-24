package main

import "fmt"

func main() {
	a := []int{}
	b := []int{1, 2, 3}
	c := a
	// a受扩容影响 内存地址发生变化 但c仍旧指向变化前的地址
	a = append(b, 1)
	fmt.Printf("address of a:%p\n", a)
	fmt.Printf("address of c:%p\n", c)
}
