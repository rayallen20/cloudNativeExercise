package main

import "fmt"

func main() {
	a := 1
	passByValue(a)
	fmt.Printf("after passByValue, a = %d\n", a)

	fmt.Printf("address of a in main:%p\n", &a)
	passByReference(&a)
	fmt.Printf("after passByReference, a = %d\n", a)
}

func passByValue(a int) {
	a += 1
	fmt.Printf("in passByValue, a = %d\n", a)
}

func passByReference(a *int) {
	fmt.Printf("address of a in func:%p\n", a)
	*a += 1
	fmt.Printf("in passByReference, a = %d\n", *a)
}
