package repository_database_interfaces

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type PodcastRepository interface {
	CreateEpisode(db *gorm.DB, episode *entities.Episode) error
	CreatePodcast(db *gorm.DB, podcast *entities.Podcast) error
	DeleteEpisode(db *gorm.DB, episodeID uint) error
	DeletePodcast(db *gorm.DB, podcastID uint) error
	ExistSubscriberByID(db *gorm.DB, podcast *entities.Podcast, userID uint) bool
	FindAllEpisodes(db *gorm.DB, offset int, pageSize int) ([]*entities.Episode, bool)
	FindAllPodcasts(db *gorm.DB, offset int, pageSize int) ([]*entities.Podcast, bool)
	FindDetailedPodcastByID(db *gorm.DB, podcastID uint) (*entities.Podcast, bool)
	FindEpisodeByID(db *gorm.DB, episodeID uint) (*entities.Episode, bool)
	FindPodcastByID(db *gorm.DB, podcastID uint) (*entities.Podcast, bool)
	FindPodcastByName(db *gorm.DB, name string) (*entities.Podcast, bool)
	FindPodcastEpisodeByName(db *gorm.DB, name string, podcastID uint) (*entities.Episode, bool)
	FindPodcastsByCategoryName(db *gorm.DB, categories []string, offset int, pageSize int) []*entities.Podcast
	FullTextSearch(db *gorm.DB, query string, offset int, pageSize int) []*entities.Podcast
	SubscribePodcast(db *gorm.DB, podcast *entities.Podcast, user *entities.User)
	UnSubscribePodcast(db *gorm.DB, podcast *entities.Podcast, user *entities.User)
	UpdateEpisode(db *gorm.DB, episode *entities.Episode) error
	UpdatePodcast(db *gorm.DB, podcast *entities.Podcast) error
	UpdatePodcastCategories(db *gorm.DB, podcastID uint, categories []entities.Category)
}
