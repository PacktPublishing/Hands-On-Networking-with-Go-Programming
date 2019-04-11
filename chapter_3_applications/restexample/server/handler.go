package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
	"github.com/gorilla/mux"
)

type CompanyStore interface {
	Get(id int64) (company company.Company, ok bool, err error)
	List() (companies []company.Company, err error)
	Upsert(c company.Company) (id int64, err error)
	Delete(id int64) (err error)
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
	log.Println(r.Method, r.URL, "CompanyHandler.List")
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
	log.Println(r.Method, r.URL, "CompanyHandler.Post")
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
	log.Println(r.Method, r.URL, "CompanyHandler.Get")
	idVar, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "missing ID", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(idVar, 10, 64)
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

func (ch *CompanyHandler) GetMany(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL, "CompanyHandler.GetMany")
	idVar, ok := mux.Vars(r)["ids"]
	if !ok {
		http.Error(w, "missing IDs", http.StatusNotFound)
		return
	}
	idStrings := strings.Split(idVar, ",")
	cs := make([]company.Company, len(idStrings))
	for i, idsv := range idStrings {
		id, err := strconv.ParseInt(idsv, 10, 64)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}
		c, ok, err := ch.CompanyStore.Get(id)
		if err != nil {
			http.Error(w, "error getting company from store", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "ID not found", http.StatusNotFound)
			return
		}
		cs[i] = c
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(cs); err != nil {
		http.Error(w, "failed to encode result", http.StatusInternalServerError)
		return
	}
	return
}

func (ch *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL, "CompanyHandler.Delete")
	idVar, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "missing ID", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(idVar, 10, 64)
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
