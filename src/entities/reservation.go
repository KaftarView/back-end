package entities

import (
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	UserID        uint                    `gorm:"not null"`
	Expiration    time.Time               `gorm:"not null"`
	Status        enums.ReservationStatus `gorm:"type:int;not null"`
	TotalPrice    float64                 `gorm:"type:decimal(10,2);not null"`
	DiscountID    *uint                   `gorm:"null"`
	Discount      *Discount               `gorm:"foreignKey:DiscountID"`
	DiscountType  *enums.DiscountType     `gorm:"type:int;null"`
	DiscountValue *float64                `gorm:"null"`
	Items         []*ReservationItem      `gorm:"foreignKey:ReservationID"`
}

type ReservationItem struct {
	gorm.Model
	ReservationID  uint    `gorm:"not null"`
	TicketID       uint    `gorm:"not null"`
	TicketName     string  `gorm:"type:varchar(50);not null"`
	TicketPrice    float64 `gorm:"type:decimal(10,2);not null"`
	TicketQuantity uint    `gorm:"not null"`
	Ticket         Ticket  `gorm:"foreignKey:TicketID"`
}
