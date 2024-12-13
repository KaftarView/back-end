package application

import (
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/exceptions"
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

func (podcastService *PodcastService) CreatePodcast(name, description string, categoryNames []string, publisherID uint) uint {
	var conflictError exceptions.ConflictError
	_, podcastExist := podcastService.podcastRepository.FindPodcastByName(name)
	if podcastExist {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Tittle,
			podcastService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
	categories := podcastService.podcastRepository.FindCategoriesByNames(categoryNames)
	commentable := podcastService.commentRepository.CreateNewCommentable()

	podcastModel := entities.Podcast{
		ID:          commentable.CID,
		Name:        name,
		Description: description,
		PublisherID: publisherID,
		Categories:  categories,
	}

	podcast := podcastService.podcastRepository.CreatePodcast(podcastModel)
	return podcast.ID
}

func (podcastService *PodcastService) SetPodcastBannerPath(bannerPath string, podcastID uint) {
	var conflictError exceptions.ConflictError
	podcast, podcastExist := podcastService.podcastRepository.FindPodcastByID(podcastID)
	if !podcastExist {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Podcast,
			podcastService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
	podcast.BannerPath = bannerPath
	podcastService.podcastRepository.UpdatePodcast(&podcast)
}

func (podcastService *PodcastService) UpdatePodcast(podcastID uint, name, description *string, Categories *[]string) entities.Podcast {
	var conflictError exceptions.ConflictError
	podcast, podcastExist := podcastService.podcastRepository.FindPodcastByID(podcastID)
	if !podcastExist {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Podcast,
			podcastService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	if name != nil {
		podcast.Name = *name
	}
	if description != nil && *description != "" {
		podcast.Description = *description
	}
	if Categories != nil {
		podcast.Categories = podcastService.podcastRepository.FindCategoriesByNames(*Categories)
	}
	podcastService.podcastRepository.UpdatePodcast(&podcast)

	return podcast
}
