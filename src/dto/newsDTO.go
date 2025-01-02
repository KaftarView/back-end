package dto

import (
	"mime/multipart"
	"time"
)

type CreateNewsRequest struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Content     string                `json:"content"`
	Content2    string                `json:"content2"`
	Banner      *multipart.FileHeader `json:"banner"`
	Banner2     *multipart.FileHeader `json:"banner2"`
	Categories  []string              `json:"categories"`
	AuthorID    uint                  `json:"authorID"`
}

type UpdateNewsRequest struct {
	ID          uint                  `json:"newsID"`
	Title       *string               `json:"title"`
	Description *string               `json:"description"`
	Content     *string               `json:"content"`
	Content2    *string               `json:"content2"`
	Banner      *multipart.FileHeader `json:"banner"`
	Banner2     *multipart.FileHeader `json:"banner2"`
	Categories  *[]string             `json:"categories"`
}

type NewsDetailsResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Content2    string    `json:"content2"`
	Banner      string    `json:"banner"`
	Banner2     string    `json:"banner2"`
	Categories  []string  `json:"categories"`
	Author      string    `json:"author"`
}
