package application_interfaces

import "first-project/src/dto"

type CommentService interface {
	CreateComment(authorID uint, commentableID uint, content string)
	DeleteComment(authorID uint, commentID uint, canModerateComment bool)
	EditComment(authorID uint, commentID uint, newContent string)
	GetPostComments(commentableID uint) []dto.CommentDetailsResponse
}
