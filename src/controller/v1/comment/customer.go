package controller_v1_comment

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type CustomerCommentController struct {
	constants      *bootstrap.Constants
	commentService application_interfaces.CommentService
}

func NewCustomerCommentController(
	constants *bootstrap.Constants,
	commentService application_interfaces.CommentService,
) *CustomerCommentController {
	return &CustomerCommentController{
		constants:      constants,
		commentService: commentService,
	}
}

func (customerCommentController *CustomerCommentController) CreateComment(c *gin.Context) {
	type createCommentParams struct {
		PostID  uint   `uri:"postID" validate:"required"`
		Content string `json:"content" validate:"required"`
	}
	param := controller.Validated[createCommentParams](c, &customerCommentController.constants.Context)
	userID, _ := c.Get(customerCommentController.constants.Context.UserID)
	customerCommentController.commentService.CreateComment(userID.(uint), param.PostID, param.Content)

	trans := controller.GetTranslator(c, customerCommentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.addComment")
	controller.Response(c, 200, message, nil)
}

func (customerCommentController *CustomerCommentController) EditComment(c *gin.Context) {
	type editCommentParams struct {
		CommentID uint   `uri:"commentID" validate:"required"`
		Content   string `json:"content" validate:"required"`
	}
	param := controller.Validated[editCommentParams](c, &customerCommentController.constants.Context)
	userID, _ := c.Get(customerCommentController.constants.Context.UserID)
	customerCommentController.commentService.EditComment(userID.(uint), param.CommentID, param.Content)

	trans := controller.GetTranslator(c, customerCommentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.editComment")
	controller.Response(c, 200, message, nil)
}

func (customerCommentController *CustomerCommentController) DeleteComment(c *gin.Context) {
	type deleteCommentParams struct {
		CommentID uint `uri:"commentID" validate:"required"`
	}
	param := controller.Validated[deleteCommentParams](c, &customerCommentController.constants.Context)
	userID, _ := c.Get(customerCommentController.constants.Context.UserID)
	customerCommentController.commentService.DeleteComment(userID.(uint), param.CommentID, false)

	trans := controller.GetTranslator(c, customerCommentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteComment")
	controller.Response(c, 200, message, nil)
}
