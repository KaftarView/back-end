package controller_v1_comment

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type AdminCommentController struct {
	constants      *bootstrap.Constants
	commentService *application.CommentService
}

func NewAdminCommentController(
	constants *bootstrap.Constants,
	commentService *application.CommentService,
) *AdminCommentController {
	return &AdminCommentController{
		constants:      constants,
		commentService: commentService,
	}
}

func (adminCommentController *AdminCommentController) DeleteComment(c *gin.Context) {
	type deleteCommentParams struct {
		CommentID uint `uri:"commentID" validate:"required"`
	}
	param := controller.Validated[deleteCommentParams](c, &adminCommentController.constants.Context)
	userID, _ := c.Get(adminCommentController.constants.Context.UserID)
	adminCommentController.commentService.DeleteComment(userID.(uint), param.CommentID, true)

	trans := controller.GetTranslator(c, adminCommentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteComment")
	controller.Response(c, 200, message, nil)
}
