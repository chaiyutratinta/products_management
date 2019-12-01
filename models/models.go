package models

//Products detail
type Products struct {
	ID       string
	Name     string
	Exp      string
	Category string
	Amount   int
	Price    int
}

type Body struct {
	Name     string `json:"name"`
	Exp      string `json:"expire_date"`
	Category string `json:"category"`
	Amount   int    `json:"amount"`
	Price    int    `json:"price"`
}
