package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type reservationRepository struct{}

func NewReservationRepository(db *gorm.DB) *reservationRepository {
	return &reservationRepository{}
}

func (reservationRepository *reservationRepository) GetExpiredReservations(db *gorm.DB) ([]*entities.Reservation, error) {
	var expiredReservations []*entities.Reservation
	err := db.
		Preload("Items").
		Where("status = ? AND expiration < ?", enums.Pending, time.Now()).
		Find(&expiredReservations).Error

	return expiredReservations, err
}

func (reservationRepository *reservationRepository) UpdateReservation(db *gorm.DB, reservation *entities.Reservation) error {
	return db.Save(reservation).Error
}
