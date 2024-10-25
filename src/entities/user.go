package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name              string
	Email             string
	Password          string
	Token             string
	Verified          bool
	Roles             []Role     `gorm:"many2many:user_roles;"`
	PreviousPasswords []Password `gorm:"foreignKey:UserID"`
}
