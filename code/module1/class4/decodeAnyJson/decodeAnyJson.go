package main

import (
	"encoding/json"
	"fmt"
)

func main() {

}

func decodeAny(humanStr string) {
	var obj interface{}
	err := json.Unmarshal([]byte(humanStr), &obj)
	if err != nil {
		return
	}
	objMap := obj.(map[string]interface{})
	for k, v := range objMap {
		switch value := v.(type) {
		case string:
			fmt.Printf("type of %s is string, value is %v\n", k, value)
		case interface{}:
			fmt.Printf("type of %s is interface{}, value is %v\n", k, value)
		default:
			fmt.Printf("type of %s is wrong, value is %v\n", k, value)
		}
	}
}
