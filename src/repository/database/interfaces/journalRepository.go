package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type JournalRepository interface {
	CreateJournal(db *gorm.DB, journal *entities.Journal) error
	DeleteJournal(db *gorm.DB, journalID uint) error
	FindAllJournals(db *gorm.DB, offset int, pageSize int) ([]*entities.Journal, bool)
	FindJournalByID(db *gorm.DB, journalID uint) (*entities.Journal, bool)
	FindJournalByName(db *gorm.DB, name string) (*entities.Journal, bool)
	FullTextSearch(db *gorm.DB, query string, offset int, pageSize int) []*entities.Journal
	UpdateJournal(db *gorm.DB, journal *entities.Journal) error
}
