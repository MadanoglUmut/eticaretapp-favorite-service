package models

import "errors"

var ErrUserNotFound error = errors.New("userID bulunamadi")

var ErrRecordNotFound error = errors.New("Kayıt Bulunamadi")

var ErrunaUthorizedAction error = errors.New("Kullanıcı Sadece Kendi Listesine Erişebilir")

var ErrUserUnauthorized error = errors.New("Kullanici Bulunamadi")
