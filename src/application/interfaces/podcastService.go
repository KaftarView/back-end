package application_interfaces

import (
	"first-project/src/dto"
	"first-project/src/entities"
	"mime/multipart"
)

type PodcastService interface {
	CreateEpisode(
		name string, description string, banner *multipart.FileHeader,
		audio *multipart.FileHeader, podcastID uint, publisherID uint) *entities.Episode
	CreatePodcast(name string, description string, categoryNames []string, banner *multipart.FileHeader, publisherID uint) *entities.Podcast
	DeleteEpisode(episodeID uint)
	DeletePodcast(podcastID uint)
	FilterPodcastsByCategory(categories []string, page int, pageSize int) []dto.PodcastDetailsResponse
	GetEpisodeDetails(episodeID uint) dto.EpisodeDetailsResponse
	GetEpisodesList(page int, pageSize int) []dto.EpisodeDetailsResponse
	GetPodcastDetails(podcastID uint) dto.PodcastDetailsResponse
	GetPodcastList(page int, pageSize int) []dto.PodcastDetailsResponse
	IsUserSubscribedPodcast(podcastID uint, userID uint) bool
	SearchEvents(query string, page int, pageSize int) []dto.PodcastDetailsResponse
	SubscribePodcast(podcastID uint, userID uint)
	UnSubscribePodcast(podcastID uint, userID uint)
	UpdateEpisode(episodeID uint, name *string, description *string, banner *multipart.FileHeader, audio *multipart.FileHeader)
	UpdatePodcast(podcastID uint, name *string, description *string, categories *[]string, banner *multipart.FileHeader)
}
