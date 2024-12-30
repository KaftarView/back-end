package repository_database_interfaces

import "first-project/src/entities"

type JournalRepository interface {
	CreateJournal(journal *entities.Journal) *entities.Journal
	DeleteJournal(journalID uint)
	FindAllJournals(offset int, pageSize int) ([]*entities.Journal, bool)
	FindJournalByID(journalID uint) (*entities.Journal, bool)
	FindJournalByName(name string) (*entities.Journal, bool)
	FullTextSearch(query string, offset int, pageSize int) []*entities.Journal
	UpdateJournal(journal *entities.Journal)
}
