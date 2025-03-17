package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FavoriteList struct {
	Id          int       `json:"list_id" gorm:"autoIncrement;column:id"`
	ListName    string    `json:"list_name" gorm:"column:listname"`
	CreatedDate time.Time `json:"created_date" gorm:"column:createddate;default:now()"`
	UserId      int       `json:"user_id" gorm:"column:userid"`
}

type FavoriteListResponse struct {
	ListId   int       `json:"list_id"`
	ListName string    `json:"list_name"`
	Items    []Product `json:"products"`
}

type CreateFavoriteList struct {
	ListName string `json:"list_name"`
	UserId   int    `json:"user_id"`
}

type UpdateFavoriteList struct {
	ListName string `json:"list_name"`
}

func (a CreateFavoriteList) Validate() error {

	return validation.ValidateStruct(&a,
		validation.Field(&a.ListName, validation.Required.Error("İsim alanı Zorunlu"), validation.Length(2, 100).Error("İsim 2-100 aralığında olmalı")))
	//validation.Field(&a.UserId, validation.Required))

}

func (a UpdateFavoriteList) Validate() error {

	return validation.ValidateStruct(&a,
		validation.Field(&a.ListName, validation.Required, validation.Length(2, 100)))

}
