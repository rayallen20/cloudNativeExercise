package main

import "fmt"

func main() {
	arr := []int{1, 2, 3}
	for _, v := range arr {
		arr = append(arr, v)
	}

	fmt.Println(arr)
}
