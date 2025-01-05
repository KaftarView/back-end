package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type PurchaseRepository interface {
	GetExpiredReservations(db *gorm.DB) ([]*entities.Reservation, error)
	GetReservationByID(db *gorm.DB, reservationID uint) (*entities.Reservation, bool)
	CreateReservation(db *gorm.DB, reservation *entities.Reservation) error
	UpdateReservation(db *gorm.DB, reservation *entities.Reservation) error
	CreateTransaction(db *gorm.DB, transaction *entities.Transaction) error
	CreateOrder(db *gorm.DB, order *entities.Order) error
	GetUserOrders(db *gorm.DB, userID uint) []*entities.Order
}
