package controller_v1_private

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	constants      *bootstrap.Constants
	commentService *application.CommentService
}

func NewCommentController(constants *bootstrap.Constants, commentService *application.CommentService) *CommentController {
	return &CommentController{
		constants:      constants,
		commentService: commentService,
	}
}

func (commentController *CommentController) GetComments(c *gin.Context) {
	type getPostCommentsParams struct {
		PostID uint `uri:"postID" validate:"required"`
	}
	param := controller.Validated[getPostCommentsParams](c, &commentController.constants.Context)
	comments := commentController.commentService.GetPostComments(param.PostID)

	controller.Response(c, 200, "", comments)
}

func (commentController *CommentController) CreateComment(c *gin.Context) {
	type createCommentParams struct {
		PostID  uint   `uri:"postID" validate:"required"`
		Content string `json:"content" validate:"required"`
	}
	param := controller.Validated[createCommentParams](c, &commentController.constants.Context)
	userID, _ := c.Get(commentController.constants.Context.UserID)
	commentController.commentService.CreateComment(userID.(uint), param.PostID, param.Content)

	trans := controller.GetTranslator(c, commentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.addComment")
	controller.Response(c, 200, message, nil)
}

func (commentController *CommentController) EditComment(c *gin.Context) {
	type editCommentParams struct {
		CommentID uint   `uri:"commentID" validate:"required"`
		Content   string `json:"content" validate:"required"`
	}
	param := controller.Validated[editCommentParams](c, &commentController.constants.Context)
	userID, _ := c.Get(commentController.constants.Context.UserID)
	commentController.commentService.EditComment(userID.(uint), param.CommentID, param.Content)

	trans := controller.GetTranslator(c, commentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.editComment")
	controller.Response(c, 200, message, nil)
}

func (commentController *CommentController) DeleteCommentByUser(c *gin.Context) {
	type deleteCommentParams struct {
		CommentID uint `uri:"commentID" validate:"required"`
	}
	param := controller.Validated[deleteCommentParams](c, &commentController.constants.Context)
	userID, _ := c.Get(commentController.constants.Context.UserID)
	commentController.commentService.DeleteComment(userID.(uint), param.CommentID, false)

	trans := controller.GetTranslator(c, commentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteComment")
	controller.Response(c, 200, message, nil)
}

func (commentController *CommentController) DeleteCommentByAdmin(c *gin.Context) {
	type deleteCommentParams struct {
		CommentID uint `uri:"commentID" validate:"required"`
	}
	param := controller.Validated[deleteCommentParams](c, &commentController.constants.Context)
	userID, _ := c.Get(commentController.constants.Context.UserID)
	commentController.commentService.DeleteComment(userID.(uint), param.CommentID, true)

	trans := controller.GetTranslator(c, commentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteComment")
	controller.Response(c, 200, message, nil)
}
