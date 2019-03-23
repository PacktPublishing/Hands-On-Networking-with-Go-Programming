package main

import (
	"sort"
	"sync"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
)

func NewInMemoryCompanyStore() *InMemoryCompanyStore {
	return &InMemoryCompanyStore{
		companies: make(map[int]company.Company),
		m:         &sync.Mutex{},
	}
}

type InMemoryCompanyStore struct {
	companies map[int]company.Company
	m         *sync.Mutex
}

func (imcs *InMemoryCompanyStore) List() (companies []company.Company, err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	// Get all keys.
	keys := make([]int, len(imcs.companies))
	var i int
	for k := range imcs.companies {
		keys[i] = k
		i++
	}
	// Sort the keys.
	sort.Ints(keys)
	// Return the sorted result.
	companies = make([]company.Company, len(keys))
	for i, k := range keys {
		companies[i] = imcs.companies[k]
	}
	return
}

func (imcs *InMemoryCompanyStore) Upsert(c company.Company) (id int, err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	if c.ID == 0 {
		c.ID = len(imcs.companies)
	}
	imcs.companies[c.ID] = c
	id = c.ID
	return
}

func (imcs *InMemoryCompanyStore) Get(id int) (c company.Company, ok bool, err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	c, ok = imcs.companies[id]
	return
}

func (imcs *InMemoryCompanyStore) Delete(id int) (err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	delete(imcs.companies, id)
	return
}
