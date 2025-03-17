package models

type Users struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Isim     string `json:"isim"`
	Soyisim  string `json:"soyisim"`
	Resim    string `json:"resim"`
}

type UserResponse struct {
	SuccesData Users `json:"SuccesData"`
}
