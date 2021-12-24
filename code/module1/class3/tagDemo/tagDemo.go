package main

import (
	"fmt"
	"reflect"
)

type MyType struct {
	Name string `json:"name"`
}

func main() {
	mt := MyType{Name: "test"}
	myType := reflect.TypeOf(mt)
	name := myType.Field(0)
	tag := name.Tag.Get("json")
	fmt.Println(tag)
}
