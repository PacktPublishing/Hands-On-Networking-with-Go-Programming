package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("backend 2")
	http.ListenAndServe(":1111", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request:", r.URL.String())
		w.Write([]byte("backend2"))
	}))
}
