package main

func main() {
	myMap := map[string]string{
		"a": "a1",
		"c": "c3",
		"d": "d4",
		"e": "e5",
	}
	element, exist := myMap["b"]
	if exist {
		println(element)
	}

	// 遍历Map
	for key, value := range myMap {
		println(key, value)
	}
}
