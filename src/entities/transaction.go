package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ReservationID uint                    `gorm:"not null"`
	Reservation   Reservation             `gorm:"foreignKey:ReservationID"`
	UserID        uint                    `gorm:"not null"`
	Amount        float64                 `gorm:"type:decimal(10,2);not null"`
	Status        enums.TransactionStatus `gorm:"type:int;not null"`
	GatewayName   string                  `gorm:"type:varchar(100);not null"`
	TrackingID    string                  `gorm:"type:varchar(100);not null"`
	ErrorMessage  *string                 `gorm:"type:varchar(200)"`
}
