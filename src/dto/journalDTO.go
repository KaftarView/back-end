package dto

import "time"

type JournalDetailsResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
	Banner      string    `json:"banner"`
	JournalFile string    `json:"journalFile"`
	Author      string    `json:"author"`
}
