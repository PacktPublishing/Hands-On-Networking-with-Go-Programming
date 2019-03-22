package company

// ID is returned when a company is posted.
type ID struct {
	ID int `json:"id"`
}

type Company struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	VATNumber string  `json:"vatNumber"`
	Address   Address `json:"address"`
}

type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	Address4 string `json:"address4"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}
