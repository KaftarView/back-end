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

func (commentController *CommentController) CreateEvent(c *gin.Context) {
	type createCommentParams struct {
		AuthorID uint   `json:"userID" validate:"required"`
		PostID   uint   `uri:"postID" validate:"required"`
		Content  string `json:"content" validate:"required"`
	}
	param := controller.Validated[createCommentParams](c, &commentController.constants.Context)
	commentController.commentService.CreateComment(param.AuthorID, param.PostID, param.Content)

	trans := controller.GetTranslator(c, commentController.constants.Context.Translator)
	message, _ := trans.T("successMessage.addComment")
	controller.Response(c, 200, message, nil)
}
