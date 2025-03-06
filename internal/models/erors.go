package models

import "errors"

var ErrUserNotFound error = errors.New("userID bulunamadi")

var ErrRecordNotFound error = errors.New("Kayıt Bulunamadi")
