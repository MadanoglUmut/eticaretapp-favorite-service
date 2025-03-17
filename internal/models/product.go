package models

type ProductResponse struct {
	SuccesData Product `json:"SuccesData"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
