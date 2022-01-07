package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
)

func main() {
	foo, bar := "foo", 1
	fmt.Println(foo, bar)

	mux := http.NewServeMux()
	// Tips: 业务逻辑的handleFunc 要挂在mux上
	mux.HandleFunc("/", rootHandler)
	// 注册pprof的handler 这些handleFunc都是内置的
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	err := http.ListenAndServe(":80", mux)
	if err != nil {
		panic(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===========Details of the http request header:=====\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s = %s\n", k, v))
	}
}
