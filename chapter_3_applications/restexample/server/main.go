package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
	"github.com/gorilla/mux"
)

func main() {
	m := mux.NewRouter()
	ch := NewCompanyHandler()
	m.Path("/companies").Methods("GET").HandlerFunc(ch.List)
	m.Path("/companies").Methods("POST").HandlerFunc(ch.Post)
	m.Path("/company/{id}").Methods("GET").HandlerFunc(ch.Get)
	m.Path("/company/{id}").Methods("POST").HandlerFunc(ch.Post)
	m.Path("/company/{id}/delete").Methods("POST").HandlerFunc(ch.Delete)
	http.ListenAndServe(":9021", m)
}

func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{
		companies: make(map[int]company.Company),
		m:         &sync.Mutex{},
	}
}

type CompanyHandler struct {
	companies map[int]company.Company
	m         *sync.Mutex
}

func (ch *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	ch.m.Lock()
	defer ch.m.Unlock()
	// Get all keys.
	keys := make([]int, len(ch.companies))
	var i int
	for k := range ch.companies {
		keys[i] = k
		i++
	}
	// Sort the keys.
	sort.Ints(keys)
	// Return the sorted result.
	result := make([]company.Company, len(keys))
	for i, k := range keys {
		result[i] = ch.companies[k]
	}
	// Encode the output.
	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		http.Error(w, "failed to encode result", http.StatusInternalServerError)
		return
	}
	return
}

func (ch *CompanyHandler) Post(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	ch.m.Lock()
	defer ch.m.Unlock()
	dec := json.NewDecoder(r.Body)
	var c company.Company
	if err := dec.Decode(&c); err != nil {
		http.Error(w, "failed to decode body", http.StatusBadRequest)
		return
	}
	if c.ID == 0 {
		c.ID = len(ch.companies)
	}
	ch.companies[c.ID] = c
	result := company.ID{ID: c.ID}
	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		http.Error(w, "failed to encode result", http.StatusInternalServerError)
		return
	}
	return
}

func (ch *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	ch.m.Lock()
	defer ch.m.Unlock()
	idVar, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "missing ID", http.StatusNotFound)
		return
	}
	id, err := strconv.Atoi(idVar)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(w)
	c, ok := ch.companies[id]
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
	}
	if err := enc.Encode(c); err != nil {
		http.Error(w, "failed to encode result", http.StatusInternalServerError)
		return
	}
	return
}

func (ch *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	ch.m.Lock()
	defer ch.m.Unlock()
	idVar, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "missing ID", http.StatusNotFound)
		return
	}
	id, err := strconv.Atoi(idVar)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	delete(ch.companies, id)
	w.Write([]byte(`{ "ok": true; }`))
}
