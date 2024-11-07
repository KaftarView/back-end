package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Type enums.PermissionType `gorm:"not null"`
}
