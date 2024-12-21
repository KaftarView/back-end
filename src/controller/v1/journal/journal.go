package controller_v1_journal

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"mime/multipart"

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
	type getJournalsListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[getJournalsListParams](c, &journalController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	journals := journalController.journalService.GetJournalsList(param.Page, param.PageSize)

	controller.Response(c, 200, "", journals)
}

func (journalController *JournalController) CreateJournal(c *gin.Context) {
	type createJournalParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Description string                `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		JournalFile *multipart.FileHeader `form:"file"`
	}
	param := controller.Validated[createJournalParams](c, &journalController.constants.Context)
	userID, _ := c.Get(journalController.constants.Context.UserID)
	journal := journalController.journalService.CreateJournal(param.Name, param.Description, param.Banner, param.JournalFile, userID.(uint))

	trans := controller.GetTranslator(c, journalController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createJournal")
	controller.Response(c, 200, message, journal.ID)
}

func (journalController *JournalController) UpdateJournal(c *gin.Context) {
	type updateJournalParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Description *string               `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		JournalFile *multipart.FileHeader `form:"file"`
		JournalID   uint                  `uri:"journalID" validate:"required"`
	}
	param := controller.Validated[updateJournalParams](c, &journalController.constants.Context)
	journalController.journalService.UpdateJournal(param.JournalID, param.Name, param.Description, param.Banner, param.JournalFile)

	trans := controller.GetTranslator(c, journalController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateJournal")
	controller.Response(c, 200, message, nil)
}

func (journalController *JournalController) DeleteJournal(c *gin.Context) {
	type deleteJournalParams struct {
		JournalID uint `uri:"journalID" validate:"required"`
	}
	param := controller.Validated[deleteJournalParams](c, &journalController.constants.Context)
	journalController.journalService.DeleteJournal(param.JournalID)

	trans := controller.GetTranslator(c, journalController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteJournal")
	controller.Response(c, 200, message, nil)
}

func (journalController *JournalController) SearchJournals(c *gin.Context) {
	// some code here ...
}
