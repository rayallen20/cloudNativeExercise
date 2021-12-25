package main

func main() {
	defer func() {
		if r := recover(); r != nil {
			println("recovered in Func Main")
		}
	}()
}
