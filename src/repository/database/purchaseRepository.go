package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type purchaseRepository struct{}

func NewPurchaseRepository(db *gorm.DB) *purchaseRepository {
	return &purchaseRepository{}
}

func (repo *purchaseRepository) GetExpiredReservations(db *gorm.DB) ([]*entities.Reservation, error) {
	var expiredReservations []*entities.Reservation
	err := db.
		Preload("Items").
		Where("status = ? AND expiration < ?", enums.Pending, time.Now()).
		Find(&expiredReservations).Error

	return expiredReservations, err
}

func (repo *purchaseRepository) GetReservationByID(db *gorm.DB, reservationID uint) (*entities.Reservation, bool) {
	var reservation entities.Reservation
	result := db.Preload("Items").First(&reservation, "id = ?", reservationID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &reservation, true
}

func (repo *purchaseRepository) CreateReservation(db *gorm.DB, reservation *entities.Reservation) error {
	return db.Create(reservation).Error
}

func (purchaseRepository *purchaseRepository) UpdateReservation(db *gorm.DB, reservation *entities.Reservation) error {
	return db.Save(reservation).Error
}

func (repo *purchaseRepository) CreateTransaction(db *gorm.DB, transaction *entities.Transaction) error {
	return db.Create(transaction).Error
}

func (repo *purchaseRepository) CreateOrder(db *gorm.DB, order *entities.Order) error {
	return db.Create(order).Error
}
