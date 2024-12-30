package application_interfaces

import (
	"first-project/src/dto"
	"first-project/src/entities"
	"mime/multipart"
)

type JournalService interface {
	CreateJournal(name string, description string, banner *multipart.FileHeader, journalFile *multipart.FileHeader, authorID uint) *entities.Journal
	DeleteJournal(journalID uint)
	GetJournalsList(page int, pageSize int) []dto.JournalDetailsResponse
	SearchJournals(query string, page int, pageSize int) []dto.JournalDetailsResponse
	UpdateJournal(journalID uint, name *string, description *string, banner *multipart.FileHeader, journalFile *multipart.FileHeader)
}
