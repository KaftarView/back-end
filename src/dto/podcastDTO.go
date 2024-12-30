package dto

import "time"

type PodcastDetailsResponse struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Banner           string    `json:"banner"`
	Publisher        string    `json:"publisher"`
	Categories       []string  `json:"categories"`
	SubscribersCount int       `json:"subscribersCount"`
}

type EpisodeDetailsResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Banner      string    `json:"banner"`
	Audio       string    `json:"audio"`
	Publisher   string    `json:"publisher"`
}
