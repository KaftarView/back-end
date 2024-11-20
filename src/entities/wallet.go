package entities

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Balance float64
	// Transactions []Transaction `gorm:"foreignKey:WalletID"`
}
