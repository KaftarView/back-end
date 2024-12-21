package repository_database

import "gorm.io/gorm"

type JournalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) *JournalRepository {
	return &JournalRepository{
		db: db,
	}
}
