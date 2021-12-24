package main

import "fmt"

func main() {
	arr := []int{1, 2, 3}
	newArr := []*int{}
	for _, v := range arr {
		newArr = append(newArr, &v)
	}

	for _, v := range newArr {
		fmt.Println(*v)
	}
}
