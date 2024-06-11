package models

type Subscription struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Price          float64  `json:"price"`
	Description    string   `json:"description"`
	Possibilities  []string `json:"possibilities"`
	DateExpiration string   `json:"date_expiration"`
}
