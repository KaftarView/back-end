package controller_v1_journal

import (
	"first-project/src/application"
	"first-project/src/bootstrap"

	"github.com/gin-gonic/gin"
)

type JournalController struct {
	constants      *bootstrap.Constants
	journalService *application.JournalService
}

func NewJournalController(
	constants *bootstrap.Constants,
	journalService *application.JournalService,
) *JournalController {
	return &JournalController{
		constants:      constants,
		journalService: journalService,
	}
}

func (journalController *JournalController) GetJournalsList(c *gin.Context) {
	// some code here ...
}

func (journalController *JournalController) CreateJournal(c *gin.Context) {
	// some code here ...
}

func (journalController *JournalController) UpdateJournal(c *gin.Context) {
	// some code here ...
}

func (journalController *JournalController) DeleteJournal(c *gin.Context) {
	// some code here ...
}

func (journalController *JournalController) SearchJournals(c *gin.Context) {
	// some code here ...
}
