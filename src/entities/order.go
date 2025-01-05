package entities

import (
	"first-project/src/enums"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID           uint                `gorm:"not null"`
	EventID          uint                `gorm:"not null"`
	ReservationID    *uint               `gorm:"null"`
	TotalPrice       float64             `gorm:"type:decimal(10,2);not null"`
	DiscountID       *uint               `gorm:"null"`
	DiscountType     *enums.DiscountType `gorm:"type:int;null"`
	DiscountValue    *float64            `gorm:"null"`
	PaymentMethod    string              `gorm:"type:varchar(50);not null"`
	PaymentReference string              `gorm:"type:varchar(100);not null"`
	User             User                `gorm:"foreignKey:UserID"`
	Reservation      *Reservation        `gorm:"foreignKey:ReservationID"`
	Items            []*OrderItem        `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID        uint    `gorm:"not null"`
	TicketID       uint    `gorm:"not null"`
	TicketName     string  `gorm:"type:varchar(50);not null"`
	TicketPrice    float64 `gorm:"type:decimal(10,2);not null"`
	TicketQuantity uint    `gorm:"not null"`
	Ticket         Ticket  `gorm:"foreignKey:TicketID"`
}
