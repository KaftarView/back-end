package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	Tag      enums.RoomType `gor:"type:int"`
	MemberID uint           `gorm:"not null"`
	Member   User           `gorm:"foreignKey:MemberID"`
	Admins   []User         `gorm:"many2many:chat_room_admins;"`
}
