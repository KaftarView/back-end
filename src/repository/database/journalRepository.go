package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type JournalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) *JournalRepository {
	return &JournalRepository{
		db: db,
	}
}

func (repo *JournalRepository) FindJournalByID(journalID uint) (*entities.Journal, bool) {
	var journal entities.Journal
	result := repo.db.First(&journal, "id = ?", journalID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &journal, true
}

func (repo *JournalRepository) FindJournalByName(name string) (*entities.Journal, bool) {
	var journal entities.Journal
	result := repo.db.First(&journal, "name = ?", name)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &journal, true
}

func (repo *JournalRepository) CreateJournal(journal *entities.Journal) *entities.Journal {
	err := repo.db.Create(journal).Error
	if err != nil {
		panic(err)
	}
	return journal
}

func (repo *JournalRepository) UpdateJournal(journal *entities.Journal) {
	err := repo.db.Save(journal).Error
	if err != nil {
		panic(err)
	}
}

func (repo *JournalRepository) DeleteJournal(journalID uint) {
	err := repo.db.Unscoped().Delete(&entities.Journal{}, journalID).Error
	if err != nil {
		panic(err)
	}
}
