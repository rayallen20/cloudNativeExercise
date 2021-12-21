package main

import "fmt"

func main() {
	Hello()
}

func Hello() {
	fmt.Println("hello")
	return
	fmt.Println("world")
}
