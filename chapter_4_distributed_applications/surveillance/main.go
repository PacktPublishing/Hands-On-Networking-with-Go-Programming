package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

func main() {
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "/Users/adrian/go/src/github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_4_distributed_applications/surveillance/"
		suffix := r.URL.Query().Get("suffix")

		// Don't do this!
		// filename := prefix + suffix
		// Do this instead.
		_, fn := path.Split(suffix)
		filename := path.Join(prefix, fn)

		f, err := os.Open(filename)
		if err != nil {
			http.Error(w, "failed to open file", http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, "failed to copy file", http.StatusInternalServerError)
			return
		}
	}))
}
