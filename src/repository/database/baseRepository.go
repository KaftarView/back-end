package repository_database

import "gorm.io/gorm"

func OrderByCreatedAtDesc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}
