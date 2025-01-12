package entities

import "gorm.io/gorm"

type ChatMessage struct {
	gorm.Model
	RoomID   uint     `gorm:"not null;index"`
	Room     ChatRoom `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SenderID uint     `gorm:"not null;index"`
	Sender   User     `gorm:"foreignKey:SenderID"`
	Content  string   `gorm:"not null;size:1000"`
}
