package main

func main() {
	DoOperation(1, increase)
	DoOperation(1, decrease)
}

func increase(a, b int) {
	println("increase result is:", a+b)
}

func decrease(a, b int) {
	println("decrease result is:", a-b)
}

func DoOperation(y int, f func(a, b int)) {
	f(y, 1)
}
