package controller_v1_journal

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralJournalController struct {
	constants      *bootstrap.Constants
	journalService application_interfaces.JournalService
}

func NewGeneralJournalController(
	constants *bootstrap.Constants,
	journalService application_interfaces.JournalService,
) *GeneralJournalController {
	return &GeneralJournalController{
		constants:      constants,
		journalService: journalService,
	}
}

func (generalJournalController *GeneralJournalController) GetJournalsList(c *gin.Context) {
	pagination := controller.GetPagination(c, &generalJournalController.constants.Context)
	journals := generalJournalController.journalService.GetJournalsList(pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", journals)
}

func (generalJournalController *GeneralJournalController) SearchJournals(c *gin.Context) {
	type searchJournalsParams struct {
		Query string `form:"query"`
	}
	param := controller.Validated[searchJournalsParams](c, &generalJournalController.constants.Context)
	pagination := controller.GetPagination(c, &generalJournalController.constants.Context)
	journals := generalJournalController.journalService.SearchJournals(param.Query, pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", journals)
}
