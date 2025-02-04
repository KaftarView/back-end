package repository_database

import (
	"first-project/src/entities"
	"strings"

	"gorm.io/gorm"
)

type JournalRepository struct{}

func NewJournalRepository() *JournalRepository {
	return &JournalRepository{}
}

func (repo *JournalRepository) FindJournalByID(db *gorm.DB, journalID uint) (*entities.Journal, bool) {
	var journal entities.Journal
	result := db.First(&journal, "id = ?", journalID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &journal, true
}

func (repo *JournalRepository) FindJournalByName(db *gorm.DB, name string) (*entities.Journal, bool) {
	var journal entities.Journal
	result := db.First(&journal, "name = ?", name)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &journal, true
}

func (repo *JournalRepository) FindAllJournals(db *gorm.DB, offset, pageSize int) ([]*entities.Journal, bool) {
	var journals []*entities.Journal
	result := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&journals)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return journals, true
}

func (repo *JournalRepository) CreateJournal(db *gorm.DB, journal *entities.Journal) error {
	return db.Create(journal).Error
}

func (repo *JournalRepository) UpdateJournal(db *gorm.DB, journal *entities.Journal) error {
	return db.Save(journal).Error
}

func (repo *JournalRepository) DeleteJournal(db *gorm.DB, journalID uint) error {
	return db.Unscoped().Delete(&entities.Journal{}, journalID).Error
}

func (repo *JournalRepository) FullTextSearch(db *gorm.DB, query string, offset, pageSize int) []*entities.Journal {
	var journals []*entities.Journal

	db.Exec(`ALTER TABLE journals ADD FULLTEXT INDEX idx_name_description (name, description)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := db.Model(&entities.Journal{}).
		Where("MATCH(name, description) AGAINST(? IN BOOLEAN MODE)", searchQuery).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&journals)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return journals
}
