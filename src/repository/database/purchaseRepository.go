package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type PurchaseRepository struct{}

func NewPurchaseRepository() *PurchaseRepository {
	return &PurchaseRepository{}
}

func (repo *PurchaseRepository) GetExpiredReservations(db *gorm.DB) ([]*entities.Reservation, error) {
	var expiredReservations []*entities.Reservation
	err := db.
		Preload("Items").
		Where("status = ? AND expiration < ?", enums.Pending, time.Now()).
		Find(&expiredReservations).Error

	return expiredReservations, err
}

func (repo *PurchaseRepository) GetReservationByID(db *gorm.DB, reservationID uint) (*entities.Reservation, bool) {
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

func (repo *PurchaseRepository) CreateReservation(db *gorm.DB, reservation *entities.Reservation) error {
	return db.Create(reservation).Error
}

func (purchaseRepository *PurchaseRepository) UpdateReservation(db *gorm.DB, reservation *entities.Reservation) error {
	return db.Save(reservation).Error
}

func (repo *PurchaseRepository) CreateTransaction(db *gorm.DB, transaction *entities.Transaction) error {
	return db.Create(transaction).Error
}

func (repo *PurchaseRepository) CreateOrder(db *gorm.DB, order *entities.Order) error {
	return db.Create(order).Error
}

func (repo *PurchaseRepository) GetUserOrders(db *gorm.DB, userID uint) []*entities.Order {
	var orders []*entities.Order

	result := OrderByCreatedAtDesc(db).Where("user_id = ?", userID).Find(&orders)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return orders
}
