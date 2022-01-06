package main

import (
	"io"
	"net/http"
)

func main() {
	// 注册handle函数
	http.HandleFunc("/shit", shit)

	// Listen and service
	err := http.ListenAndServe("127.0.0.1:8085", nil)
	if err != nil {
		panic(err)
	}
}

// handle函数
func shit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "fuck")
}
