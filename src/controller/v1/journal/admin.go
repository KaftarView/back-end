package controller_v1_journal

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type AdminJournalController struct {
	constants      *bootstrap.Constants
	journalService application_interfaces.JournalService
}

func NewAdminJournalController(
	constants *bootstrap.Constants,
	journalService application_interfaces.JournalService,
) *AdminJournalController {
	return &AdminJournalController{
		constants:      constants,
		journalService: journalService,
	}
}

func (adminJournalController *AdminJournalController) CreateJournal(c *gin.Context) {
	type createJournalParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Description string                `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		JournalFile *multipart.FileHeader `form:"file"`
	}
	param := controller.Validated[createJournalParams](c, &adminJournalController.constants.Context)
	userID, _ := c.Get(adminJournalController.constants.Context.UserID)
	journal := adminJournalController.journalService.CreateJournal(param.Name, param.Description, param.Banner, param.JournalFile, userID.(uint))

	trans := controller.GetTranslator(c, adminJournalController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createJournal")
	controller.Response(c, 200, message, journal.ID)
}

func (adminJournalController *AdminJournalController) UpdateJournal(c *gin.Context) {
	type updateJournalParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Description *string               `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		JournalFile *multipart.FileHeader `form:"file"`
		JournalID   uint                  `uri:"journalID" validate:"required"`
	}
	param := controller.Validated[updateJournalParams](c, &adminJournalController.constants.Context)
	adminJournalController.journalService.UpdateJournal(param.JournalID, param.Name, param.Description, param.Banner, param.JournalFile)

	trans := controller.GetTranslator(c, adminJournalController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateJournal")
	controller.Response(c, 200, message, nil)
}

func (adminJournalController *AdminJournalController) DeleteJournal(c *gin.Context) {
	type deleteJournalParams struct {
		JournalID uint `uri:"journalID" validate:"required"`
	}
	param := controller.Validated[deleteJournalParams](c, &adminJournalController.constants.Context)
	adminJournalController.journalService.DeleteJournal(param.JournalID)

	trans := controller.GetTranslator(c, adminJournalController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteJournal")
	controller.Response(c, 200, message, nil)
}
