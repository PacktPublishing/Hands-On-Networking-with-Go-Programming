package main

import (
	"fmt"
	"net/http"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/server"
	"github.com/gorilla/mux"
)

var hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
})

func main() {
	fmt.Println("REST company server started (port 9021)")
	m := mux.NewRouter()
	store := server.NewInMemoryCompanyStore()
	ch := server.NewCompanyHandler(store)
	m.Path("/companies").Methods("GET").HandlerFunc(ch.List)
	m.Path("/companies").Methods("POST").HandlerFunc(ch.Post)
	m.Path("/company/{id}").Methods("GET").HandlerFunc(ch.Get)
	m.Path("/company?ids={ids}").Methods("GET").HandlerFunc(ch.GetMany)
	m.Path("/company/{id}").Methods("POST").HandlerFunc(ch.Post)
	m.Path("/company/{id}/delete").Methods("POST").HandlerFunc(ch.Delete)
	http.ListenAndServe(":9021", m)
}
