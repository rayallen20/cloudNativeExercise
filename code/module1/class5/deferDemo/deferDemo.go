package main

import "fmt"

func main() {
	defer fmt.Printf("1\n")
	defer fmt.Printf("2\n")
	defer fmt.Printf("3\n")
}
