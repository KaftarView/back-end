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

func (podcastService *PodcastService) CreatePodcast(name, description string, categoryNames []string, publisherID uint) *entities.Podcast {
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

	podcast := podcastService.podcastRepository.CreatePodcast(&podcastModel)
	return podcast
}

func (podcastService *PodcastService) SetPodcastBannerPath(bannerPath string, podcast *entities.Podcast) {
	podcast.BannerPath = bannerPath
	podcastService.podcastRepository.UpdatePodcast(podcast)
}

func (podcastService *PodcastService) UpdatePodcast(podcastID uint, name, description *string, Categories *[]string) *entities.Podcast {
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
	podcastService.podcastRepository.UpdatePodcast(podcast)

	return podcast
}

func (podcastService *PodcastService) FindEpisodeByID(episodeID uint) *entities.Episode {
	var notFoundError exceptions.NotFoundError
	episode, episodeExist := podcastService.podcastRepository.FindEpisodeByID(episodeID)
	if !episodeExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Episode
		panic(notFoundError)
	}
	return episode
}

func (podcastService *PodcastService) CreateEpisode(name, description string, podcastID, publisherID uint) *entities.Episode {
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
	episodeModel := entities.Episode{
		Name:        name,
		Description: description,
		PublisherID: publisherID,
		PodcastID:   podcastID,
	}
	episode := podcastService.podcastRepository.CreateEpisode(&episodeModel)
	return episode
}

func (podcastService *PodcastService) SetEpisodeBannerPath(bannerPath string, episode *entities.Episode) {
	episode.BannerPath = bannerPath
	podcastService.podcastRepository.UpdateEpisode(episode)
}

func (podcastService *PodcastService) SetEpisodeAudioPath(audioPath string, episode *entities.Episode) {
	episode.AudioPath = audioPath
	podcastService.podcastRepository.UpdateEpisode(episode)
}

func (podcastService *PodcastService) UpdateEpisode(episodeID uint, name, description *string) *entities.Episode {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	episode, episodeExist := podcastService.podcastRepository.FindEpisodeByID(episodeID)
	if !episodeExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Episode
		panic(notFoundError)
	}

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

	podcastService.podcastRepository.UpdateEpisode(episode)
	return episode
}

func (podcastService *PodcastService) DeleteEpisode(episodeID uint) {
	var notFoundError exceptions.NotFoundError
	_, episodeExist := podcastService.podcastRepository.FindEpisodeByID(episodeID)
	if !episodeExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Episode
		panic(notFoundError)
	}
	podcastService.podcastRepository.DeleteEpisodeByID(episodeID)
}
