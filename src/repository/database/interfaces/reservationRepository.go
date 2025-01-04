package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	GetExpiredReservations(db *gorm.DB) ([]*entities.Reservation, error)
	UpdateReservation(db *gorm.DB, reservation *entities.Reservation) error
}
