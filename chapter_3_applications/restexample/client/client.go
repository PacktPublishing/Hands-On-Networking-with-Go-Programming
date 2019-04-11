package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
)

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

func (cc *CompanyClient) Get(id int64) (company company.Company, err error) {
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

func (cc *CompanyClient) GetMany(ids []int64) (companies []company.Company, err error) {
	var q url.Values
	for _, id := range ids {
		q.Add("ids", strconv.FormatInt(id, 10))
	}
	resp, err := cc.client.Get(cc.endpoint + fmt.Sprintf("/company?ids=%s", q.Encode()))
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
	bdy, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bdy, &id)
	if err != nil {
		err = fmt.Errorf("failed to decode response body '%v': %v", string(bdy), err)
	}
	return
}

func errorFromResponse(r *http.Response) error {
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("URL %s returned a 404", r.Request.URL.String())
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("URL %s returned unexpected status %d", r.Request.URL.String(), r.StatusCode)
	}
	return nil
}
