package main

import "fmt"

func main() {
	words := []string{"foo", "bar", "baz"}
	for _, word := range words {
		go func() {
			fmt.Println(word)
		}()
	}
}
