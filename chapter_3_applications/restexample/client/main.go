package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
)

func main() {
	cc := NewCompanyClient("http://localhost:9021")

	fmt.Println("Listing companies")
	companies, err := cc.List()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range companies {
		fmt.Println(c)
	}
	fmt.Println()

	fmt.Println("Creating new company")
	c := company.Company{
		Name: "Test Company 1",
		Address: company.Address{
			Address1: "Address 1",
			Address2: "Address 2",
			Address3: "Address 3",
			Address4: "Address 4",
			Postcode: "Postcode",
		},
		VATNumber: "12345",
	}
	id, err := cc.Post(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Created ID: %d\n", id)
	fmt.Println()

	fmt.Println("Getting company")
	c, err = cc.Get(id.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println()

	fmt.Println("Listing companies")
	companies, err = cc.List()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range companies {
		fmt.Println(c)
	}
	fmt.Println()
}

func NewCompanyClient(endpoint string) *CompanyClient {
	return &CompanyClient{
		endpoint: endpoint,
		client:   http.Client{},
	}
}

type CompanyClient struct {
	endpoint string
	client   http.Client
}

func (cc *CompanyClient) List() (companies []company.Company, err error) {
	resp, err := cc.client.Get(cc.endpoint + "/companies")
	if err != nil {
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		return
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&companies)
	return
}

func (cc *CompanyClient) Get(id int) (company company.Company, err error) {
	resp, err := cc.client.Get(cc.endpoint + fmt.Sprintf("/company/%d", id))
	if err != nil {
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		return
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&company)
	return
}

func (cc *CompanyClient) Post(c company.Company) (id company.ID, err error) {
	url := cc.endpoint + "/companies"
	if c.ID > 0 {
		url = cc.endpoint + fmt.Sprintf("/company/%d", c.ID)
	}
	b, err := json.Marshal(c)
	if err != nil {
		return
	}
	resp, err := cc.client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		return
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&id)
	return
}

func errorFromResponse(r *http.Response) error {
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("URL %s returned a 404", r.Request.URL.String())
	}
	return nil
}
