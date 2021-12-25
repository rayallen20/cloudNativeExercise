package main

import (
	"fmt"
	"reflect"
)

func main() {
	myMap := make(map[string]string, 10)
	myMap["a"] = "b"
	t := reflect.TypeOf(myMap)
	fmt.Printf("type:%s\n", t)
	v := reflect.ValueOf(myMap)
	fmt.Printf("value:%v\n", v)
}
