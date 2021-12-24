package main

import (
	"errors"
	"fmt"
)

func main() {
	myArray := [5]int{1, 2, 3, 4, 5}
	mySlice := myArray[1:3] // 取左闭右开区间
	fmt.Printf("%+v\n", mySlice)
	fullSlice := myArray[:] // 取数组全部元素
	remove3rdItem, _ := deleteItem(fullSlice, 2)
	fmt.Printf("remove3rdItem %+v\n", remove3rdItem)
}

// deleteItem 根据索引从数组中删除1个元素
func deleteItem(slice []int, index int) ([]int, error) {
	if index >= len(slice) {
		return nil, errors.New("illegal index")
	}
	return append(slice[:index], slice[index+1:]...), nil
}
