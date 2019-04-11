package schema

type Order struct {
	ID        string `json:"id"`
	Items     []Item `json:"items"`
	CompanyID int64  `json:"companyId"`
}
