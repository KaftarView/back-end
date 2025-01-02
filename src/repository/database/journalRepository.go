package repository_database

import (
	"first-project/src/entities"
	"strings"

	"gorm.io/gorm"
)

type journalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) *journalRepository {
	return &journalRepository{
		db: db,
	}
}

func (repo *journalRepository) FindJournalByID(journalID uint) (*entities.Journal, bool) {
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

func (repo *journalRepository) FindJournalByName(name string) (*entities.Journal, bool) {
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

func (repo *journalRepository) FindAllJournals(offset, pageSize int) ([]*entities.Journal, bool) {
	var journals []*entities.Journal
	result := repo.db.Offset(offset).Limit(pageSize).Find(&journals)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return journals, true
}

func (repo *journalRepository) CreateJournal(journal *entities.Journal) *entities.Journal {
	err := repo.db.Create(journal).Error
	if err != nil {
		panic(err)
	}
	return journal
}

func (repo *journalRepository) UpdateJournal(journal *entities.Journal) {
	err := repo.db.Save(journal).Error
	if err != nil {
		panic(err)
	}
}

func (repo *journalRepository) DeleteJournal(journalID uint) {
	err := repo.db.Unscoped().Delete(&entities.Journal{}, journalID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *journalRepository) FullTextSearch(query string, offset, pageSize int) []*entities.Journal {
	var journals []*entities.Journal

	repo.db.Exec(`ALTER TABLE journals ADD FULLTEXT INDEX idx_name_description (name, description)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := repo.db.Model(&entities.Journal{}).
		Where("MATCH(name, description) AGAINST(? IN BOOLEAN MODE)", searchQuery).
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
