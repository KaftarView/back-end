package application

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	"fmt"
	"mime/multipart"
)

type PodcastService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3service
	podcastRepository *repository_database.PodcastRepository
	commentRepository *repository_database.CommentRepository
}

func NewPodcastService(
	constants *bootstrap.Constants,
	awsS3Service *application_aws.S3service,
	podcastRepository *repository_database.PodcastRepository,
	commentRepository *repository_database.CommentRepository,
) *PodcastService {
	return &PodcastService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		podcastRepository: podcastRepository,
		commentRepository: commentRepository,
	}
}

func (podcastService *PodcastService) CreatePodcast(name, description string, categoryNames []string, banner *multipart.FileHeader, publisherID uint) *entities.Podcast {
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

	tx := podcastService.podcastRepository.BeginTransaction()

	defer tx.Rollback()

	podcast := &entities.Podcast{
		ID:          commentable.CID,
		Name:        name,
		Description: description,
		PublisherID: publisherID,
		Categories:  categories,
	}

	podcastService.podcastRepository.CreatePodcast(tx, podcast)

	if banner != nil {
		bannerPath := fmt.Sprintf("banners/podcasts/%d/images/%s", podcast.ID, banner.Filename)
		podcastService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, banner)
		podcast.BannerPath = bannerPath
		podcastService.podcastRepository.UpdatePodcast(tx, podcast)
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
	return podcast
}

func (podcastService *PodcastService) UpdatePodcast(podcastID uint, name, description *string, Categories *[]string, banner *multipart.FileHeader) {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	podcast, podcastExist := podcastService.podcastRepository.FindPodcastByID(podcastID)
	if !podcastExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Podcast
		panic(notFoundError)
	}

	tx := podcastService.podcastRepository.BeginTransaction()

	defer tx.Rollback()

	if name != nil {
		_, podcastExist := podcastService.podcastRepository.FindPodcastByName(*name)
		if podcastExist {
			conflictError.AppendError(
				podcastService.constants.ErrorField.Tittle,
				podcastService.constants.ErrorTag.AlreadyExist)
			panic(conflictError)
		}
		podcast.Name = *name
	}
	if description != nil && *description != "" {
		podcast.Description = *description
	}
	if Categories != nil {
		podcast.Categories = podcastService.podcastRepository.FindCategoriesByNames(*Categories)
	}
	if banner != nil {
		if podcast.BannerPath != "" {
			podcastService.awsS3Service.DeleteObject(enums.BannersBucket, podcast.BannerPath)
		}
		objectPath := fmt.Sprintf("banners/podcasts/%d/images/%s", podcastID, banner.Filename)
		podcastService.awsS3Service.UploadObject(enums.BannersBucket, objectPath, banner)
		podcast.BannerPath = objectPath
	}
	podcastService.podcastRepository.UpdatePodcast(tx, podcast)

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
}

func (podcastService *PodcastService) CreateEpisode(name, description string, banner, audio *multipart.FileHeader, podcastID, publisherID uint) *entities.Episode {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	_, podcastExist := podcastService.podcastRepository.FindPodcastByID(podcastID)
	if !podcastExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Podcast
		panic(notFoundError)
	}
	_, episodeExist := podcastService.podcastRepository.FindEpisodeByName(name)
	if episodeExist {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Tittle,
			podcastService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	tx := podcastService.podcastRepository.BeginTransaction()

	defer tx.Rollback()

	episode := &entities.Episode{
		Name:        name,
		Description: description,
		PublisherID: publisherID,
		PodcastID:   podcastID,
	}
	podcastService.podcastRepository.CreateEpisode(tx, episode)

	if banner != nil {
		bannerPath := fmt.Sprintf("banners/podcasts/%d/episodes/%d/images/%s", podcastID, episode.ID, banner.Filename)
		podcastService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, banner)
		episode.BannerPath = bannerPath
	}

	audioPath := fmt.Sprintf("media/podcasts/%d/episodes/%d/audio/%s", podcastID, episode.ID, audio.Filename)
	podcastService.awsS3Service.UploadObject(enums.PodcastsBucket, audioPath, audio)
	episode.AudioPath = audioPath
	podcastService.podcastRepository.UpdateEpisode(tx, episode)

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}

	return episode
}

func (podcastService *PodcastService) UpdateEpisode(episodeID uint, name, description *string, banner, audio *multipart.FileHeader) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	episode, episodeExist := podcastService.podcastRepository.FindEpisodeByID(episodeID)
	if !episodeExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Episode
		panic(notFoundError)
	}

	tx := podcastService.podcastRepository.BeginTransaction()

	defer tx.Rollback()

	if name != nil {
		_, episodeExist := podcastService.podcastRepository.FindEpisodeByName(*name)
		if episodeExist {
			conflictError.AppendError(
				podcastService.constants.ErrorField.Tittle,
				podcastService.constants.ErrorTag.AlreadyExist)
			panic(conflictError)
		}
		episode.Name = *name
	}

	if description != nil && *description != "" {
		episode.Description = *description
	}

	if banner != nil {
		if episode.BannerPath != "" {
			podcastService.awsS3Service.DeleteObject(enums.BannersBucket, episode.BannerPath)
		}
		bannerPath := fmt.Sprintf("banners/podcasts/%d/episodes/%d/images/%s", episode.PodcastID, episodeID, banner.Filename)
		podcastService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, banner)
		episode.BannerPath = bannerPath
	}

	if audio != nil {
		podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, episode.AudioPath)
		audioPath := fmt.Sprintf("media/podcasts/%d/episodes/%d/audio/%s", episode.PodcastID, episode.ID, audio.Filename)
		podcastService.awsS3Service.UploadObject(enums.PodcastsBucket, audioPath, audio)
		episode.AudioPath = audioPath
	}

	podcastService.podcastRepository.UpdateEpisode(tx, episode)

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
}

func (podcastService *PodcastService) DeleteEpisode(episodeID uint) {
	var notFoundError exceptions.NotFoundError
	episode, episodeExist := podcastService.podcastRepository.FindEpisodeByID(episodeID)
	if !episodeExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Episode
		panic(notFoundError)
	}

	tx := podcastService.podcastRepository.BeginTransaction()

	defer tx.Rollback()

	podcastService.podcastRepository.DeleteEpisodeByID(tx, episodeID)

	if episode.BannerPath != "" {
		podcastService.awsS3Service.DeleteObject(enums.BannersBucket, episode.BannerPath)
	}

	if episode.AudioPath != "" {
		podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, episode.AudioPath)
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
}
