package model

// MicroNumber
type MicroNumber struct {
	Id         int     `json:"id"`
	CurrencyId int     `json:"currency_id"`
	Number     float64 `json:"number"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
