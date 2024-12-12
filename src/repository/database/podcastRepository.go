package repository_database

import "gorm.io/gorm"

type PodcastRepository struct {
	db *gorm.DB
}

func NewPodcastRepository(db *gorm.DB) *PodcastRepository {
	return &PodcastRepository{
		db: db,
	}
}
