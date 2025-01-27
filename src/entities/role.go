package entities

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Type        string       `gorm:"not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
