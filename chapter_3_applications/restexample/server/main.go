package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

var hf = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
})

func main() {
	m := mux.NewRouter()
	store := NewInMemoryCompanyStore()
	ch := NewCompanyHandler(store)
	m.Path("/companies").Methods("GET").HandlerFunc(ch.List)
	m.Path("/companies").Methods("POST").HandlerFunc(ch.Post)
	m.Path("/company/{id}").Methods("GET").HandlerFunc(ch.Get)
	m.Path("/company/{id}").Methods("POST").HandlerFunc(ch.Post)
	m.Path("/company/{id}/delete").Methods("POST").HandlerFunc(ch.Delete)
	http.ListenAndServe(":9021", m)
}
