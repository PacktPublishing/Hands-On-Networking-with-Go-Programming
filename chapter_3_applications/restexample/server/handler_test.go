package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
)

type TestStore struct {
	GetFn    func(id int64) (company company.Company, ok bool, err error)
	ListFn   func() (companies []company.Company, err error)
	UpsertFn func(c company.Company) (id int64, err error)
	DeleteFn func(id int64) (err error)
}

func (ts TestStore) Get(id int64) (company company.Company, ok bool, err error) {
	return ts.GetFn(id)
}
func (ts TestStore) List() (companies []company.Company, err error) {
	return ts.ListFn()
}
func (ts TestStore) Upsert(c company.Company) (id int64, err error) {
	return ts.UpsertFn(c)
}
func (ts TestStore) Delete(id int64) (err error) {
	return ts.DeleteFn(id)
}

func TestList(t *testing.T) {
	tests := []struct {
		name               string
		store              CompanyStore
		r                  *http.Request
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "a store error results in an internal server error",
			store: TestStore{
				ListFn: func() (companies []company.Company, err error) {
					return nil, fmt.Errorf("database error")
				},
			},
			r:                  httptest.NewRequest("GET", "/List", nil),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "failed to get company list\n",
		},
		{
			name: "the list is returned in JSON format",
			store: TestStore{
				ListFn: func() (companies []company.Company, err error) {
					companies = []company.Company{
						company.Company{
							ID:   1,
							Name: "Test Company 1",
							Address: company.Address{
								Address1: "Address 1",
								Postcode: "Postcode",
							},
						},
					}
					return
				},
			},
			r:                  httptest.NewRequest("GET", "/List", nil),
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[{"id":1,"name":"Test Company 1","vatNumber":"","address":{"address1":"Address 1","address2":"","address3":"","address4":"","postcode":"Postcode","country":""}}]` + "\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewCompanyHandler(test.store)
			w := httptest.NewRecorder()
			h.List(w, test.r)
			actualBody, err := ioutil.ReadAll(w.Result().Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
			}
			if string(actualBody) != test.expectedBody {
				t.Errorf("expected body:\n%s\n\nactual body:\n%s\n\n", test.expectedBody, string(actualBody))
			}
			if w.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("expected status %d, got %d", test.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}
}
