package entities

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	FinalPrice    float64
	Quantity      uint
	PurchasableID uint        `gorm:"not null"`
	Purchasable   Purchasable `gorm:"foreignKey:PurchasableID"`
	// WalletID uint   `gorm:"index"`
	// Wallet   Wallet `gorm:"-"`
}
