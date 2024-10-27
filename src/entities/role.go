package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Type enums.RoleType `gorm:"not null"`
}
