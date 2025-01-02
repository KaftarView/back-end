package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type PodcastRepository interface {
	BeginTransaction() *gorm.DB
	CreateEpisode(tx *gorm.DB, episode *entities.Episode) *entities.Episode
	CreatePodcast(tx *gorm.DB, podcast *entities.Podcast) *entities.Podcast
	DeleteEpisodeByID(tx *gorm.DB, episodeID uint)
	DeletePodcast(podcastID uint)
	ExistSubscriberByID(podcast *entities.Podcast, userID uint) bool
	FindAllEpisodes(offset int, pageSize int) ([]*entities.Episode, bool)
	FindAllPodcasts(offset int, pageSize int) ([]*entities.Podcast, bool)
	FindDetailedPodcastByID(podcastID uint) (*entities.Podcast, bool)
	FindEpisodeByID(episodeID uint) (*entities.Episode, bool)
	FindPodcastByID(podcastID uint) (*entities.Podcast, bool)
	FindPodcastByName(name string) (*entities.Podcast, bool)
	FindPodcastEpisodeByName(name string, podcastID uint) (*entities.Episode, bool)
	FindPodcastsByCategoryName(categories []string, offset int, pageSize int) []*entities.Podcast
	FullTextSearch(query string, offset int, pageSize int) []*entities.Podcast
	SubscribePodcast(podcast *entities.Podcast, user *entities.User)
	UnSubscribePodcast(podcast *entities.Podcast, user *entities.User)
	UpdateEpisode(tx *gorm.DB, episode *entities.Episode)
	UpdatePodcast(tx *gorm.DB, podcast *entities.Podcast)
	UpdatePodcastCategories(podcastID uint, categories []entities.Category)
}
