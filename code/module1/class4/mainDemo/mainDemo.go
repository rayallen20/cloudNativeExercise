package main

import "os"

func main() {
	args := os.Args
	if len(args) != 0 {
		println("Do not accept any argument")
		os.Exit(1)
	}
	println("Hello world")
}
