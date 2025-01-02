package controller_v1_comment

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralCommentController struct {
	constants      *bootstrap.Constants
	commentService application_interfaces.CommentService
}

func NewGeneralCommentController(
	constants *bootstrap.Constants,
	commentService application_interfaces.CommentService,
) *GeneralCommentController {
	return &GeneralCommentController{
		constants:      constants,
		commentService: commentService,
	}
}

func (generalCommentController *GeneralCommentController) GetComments(c *gin.Context) {
	type getPostCommentsParams struct {
		PostID uint `uri:"postID" validate:"required"`
	}
	param := controller.Validated[getPostCommentsParams](c, &generalCommentController.constants.Context)
	comments := generalCommentController.commentService.GetPostComments(param.PostID)

	controller.Response(c, 200, "", comments)
}
