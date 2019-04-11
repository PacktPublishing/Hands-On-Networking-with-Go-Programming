package server

import (
	"sort"
	"sync"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
)

func NewInMemoryCompanyStore() *InMemoryCompanyStore {
	return &InMemoryCompanyStore{
		companies: make(map[int64]company.Company),
		m:         &sync.Mutex{},
	}
}

type InMemoryCompanyStore struct {
	companies map[int64]company.Company
	m         *sync.Mutex
}

func (imcs *InMemoryCompanyStore) List() (companies []company.Company, err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	// Get all keys.
	keys := make([]int64, len(imcs.companies))
	var i int
	for k := range imcs.companies {
		keys[i] = k
		i++
	}
	// Sort the keys.
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	// Return the sorted result.
	companies = make([]company.Company, len(keys))
	for i, k := range keys {
		companies[i] = imcs.companies[k]
	}
	return
}

func (imcs *InMemoryCompanyStore) Upsert(c company.Company) (id int64, err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	if c.ID == 0 {
		c.ID = int64(len(imcs.companies))
	}
	imcs.companies[c.ID] = c
	id = c.ID
	return
}

func (imcs *InMemoryCompanyStore) Get(id int64) (c company.Company, ok bool, err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	c, ok = imcs.companies[id]
	return
}

func (imcs *InMemoryCompanyStore) Delete(id int64) (err error) {
	imcs.m.Lock()
	defer imcs.m.Unlock()
	delete(imcs.companies, id)
	return
}
