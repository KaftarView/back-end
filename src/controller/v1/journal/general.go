package controller_v1_journal

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralJournalController struct {
	constants      *bootstrap.Constants
	journalService *application.JournalService
}

func NewGeneralJournalController(
	constants *bootstrap.Constants,
	journalService *application.JournalService,
) *GeneralJournalController {
	return &GeneralJournalController{
		constants:      constants,
		journalService: journalService,
	}
}

func (generalJournalController *GeneralJournalController) GetJournalsList(c *gin.Context) {
	type getJournalsListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[getJournalsListParams](c, &generalJournalController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	journals := generalJournalController.journalService.GetJournalsList(param.Page, param.PageSize)

	controller.Response(c, 200, "", journals)
}

func (generalJournalController *GeneralJournalController) SearchJournals(c *gin.Context) {
	type searchJournalsParams struct {
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	param := controller.Validated[searchJournalsParams](c, &generalJournalController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	journals := generalJournalController.journalService.SearchJournals(param.Query, param.Page, param.PageSize)

	controller.Response(c, 200, "", journals)
}
