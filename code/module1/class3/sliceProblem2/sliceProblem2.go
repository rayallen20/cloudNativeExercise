package main

import "fmt"

func main() {
	mySlice := []int{10, 20, 30, 40, 50}
	for key, value := range mySlice {
		fmt.Printf("address of value: %p\n", &value)
		fmt.Printf("address of mySlice[%d]: %p\n", key, &mySlice[key])
	}
}
