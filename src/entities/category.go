package entities

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"type:text" `
	Description string `gorm:"type:text" `
}
