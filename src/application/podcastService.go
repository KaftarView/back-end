package application

import (
	"first-project/src/bootstrap"
	repository_database "first-project/src/repository/database"
)

type PodcastService struct {
	constants         *bootstrap.Constants
	podcastRepository *repository_database.PodcastRepository
	commentRepository *repository_database.CommentRepository
}

func NewPodcastService(
	constants *bootstrap.Constants,
	podcastRepository *repository_database.PodcastRepository,
	commentRepository *repository_database.CommentRepository,
) *PodcastService {
	return &PodcastService{
		constants:         constants,
		podcastRepository: podcastRepository,
		commentRepository: commentRepository,
	}
}
