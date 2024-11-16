package entities

import "gorm.io/gorm"

type Password struct {
	gorm.Model
	Password string
	UserID   uint
}
