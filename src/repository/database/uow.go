package repository_database

import "gorm.io/gorm"

func ExecuteInTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
