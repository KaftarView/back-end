package entities

type Purchasable struct {
	ID           uint          `gorm:"primaryKey"`
	Transactions []Transaction `gorm:"foreignKey:PurchasableID"`
}
