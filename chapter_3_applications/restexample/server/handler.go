package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
	"github.com/gorilla/mux"
)

type CompanyStore interface {
	Get(id int) (company company.Company, ok bool, err error)
	List() (companies []company.Company, err error)
	Upsert(c company.Company) (id int, err error)
	Delete(id int) (err error)
}

func NewCompanyHandler(cs CompanyStore) *CompanyHandler {
	return &CompanyHandler{
		CompanyStore: cs,
	}
}

type CompanyHandler struct {
	CompanyStore CompanyStore
}

func (ch *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	result, err := ch.CompanyStore.List()
	if err != nil {
		http.Error(w, "failed to get company list", http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		http.Error(w, "failed to encode result", http.StatusInternalServerError)
		return
	}
	return
}

func (ch *CompanyHandler) Post(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	dec := json.NewDecoder(r.Body)
	var c company.Company
	if err := dec.Decode(&c); err != nil {
		http.Error(w, "failed to decode body", http.StatusBadRequest)
		return
	}
	id, err := ch.CompanyStore.Upsert(c)
	if err != nil {
		http.Error(w, "failed to upsert", http.StatusInternalServerError)
		return
	}
	result := company.ID{ID: id}
	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		http.Error(w, "failed to encode result", http.StatusInternalServerError)
		return
	}
	return
}

func (ch *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
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
	c, ok, err := ch.CompanyStore.Get(id)
	if err != nil {
		http.Error(w, "error getting company from store", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}
	if err := enc.Encode(c); err != nil {
		http.Error(w, "failed to encode result", http.StatusInternalServerError)
		return
	}
	return
}

func (ch *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
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
	err = ch.CompanyStore.Delete(id)
	if err != nil {
		http.Error(w, "error deleting company from store", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{ "ok": true; }`))
}
