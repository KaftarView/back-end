package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Type        enums.RoleType `gorm:"not null"`
	Permissions []Permission   `gorm:"many2many:role_permissions;"`
}
