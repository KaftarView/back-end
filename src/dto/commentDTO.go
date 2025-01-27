package dto

type CommentDetailsResponse struct {
	ID          uint   `json:"id"`
	Content     string `json:"content"`
	IsModerated bool   `json:"isModerated"`
	AuthorName  string `json:"authorName"`
}
