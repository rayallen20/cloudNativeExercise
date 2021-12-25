package main

import (
	"fmt"
	"reflect"
)

type T struct {
	A string
}

func (t *T) String() string {
	return t.A + "1"
}

func (t T) Change(p *T) string {
	fmt.Printf("p.A = %s\n", p.A)
	t.A = "fuck"
	return t.A
}

func main() {
	myStruct := T{A: "a"}
	v1 := reflect.ValueOf(myStruct)
	for i := 0; i < v1.NumField(); i++ {
		fmt.Printf("Field %d:%v\n", i, v1.Field(i))
	}

	for i := 0; i < v1.NumMethod(); i++ {
		// Tip: 接收者为*T的方法 不属于v1
		fmt.Printf("Method %d:%v\n", i, v1.Method(i))
	}

	// 反射调用传参示例
	param := &T{A: "c"}
	paramCall := reflect.ValueOf(param)
	result := v1.Method(0).Call([]reflect.Value{paramCall})
	fmt.Printf("%v\n", result)
}
