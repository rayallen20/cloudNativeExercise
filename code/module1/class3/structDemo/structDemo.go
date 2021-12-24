package main

type MyType struct {
	Name string
}

func printMyType(t *MyType) {
	println(t.Name)
}

func main() {
	t := MyType{Name: "test"}
	printMyType(&t)
}
