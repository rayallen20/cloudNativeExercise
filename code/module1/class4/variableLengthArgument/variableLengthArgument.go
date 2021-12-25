package main

import "fmt"

func main() {
	b := 1
	c := 2
	d := 3
	sum := accumulate(b, c, d)
	fmt.Println(sum)
}

func accumulate(a ...int) (sum int) {
	for _, v := range a {
		sum += v
	}
	return
}
