package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FavoriteItem struct {
	ItemId      int       `json:"item_id" gorm:"column:itemid"`
	ListId      int       `json:"list_id" gorm:"column:listid"`
	CreatedDate time.Time `json:"created_date" gorm:"column:createddate;default:now()"`
}

type CreateFavoriteItem struct {
	ItemId int `json:"item_id"`
	ListId int `json:"list_id"`
}

type UpdateFavoriteItem struct {
	ListId int `json:"list_id"`
}

func (a CreateFavoriteItem) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ItemId, validation.Required),
		validation.Field(&a.ListId, validation.Required))
}

func (a UpdateFavoriteItem) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ListId, validation.Required))
}
