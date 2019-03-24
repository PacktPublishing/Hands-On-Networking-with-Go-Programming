package main

import "net/http"

func main() {
	okh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServeTLS(":443", "cert.crt", "cert.key", okh)
}
