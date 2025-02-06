package application

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type PodcastService struct {
	constants         *bootstrap.Constants
	awsS3Service      application_interfaces.S3Service
	categoryService   application_interfaces.CategoryService
	podcastRepository repository_database_interfaces.PodcastRepository
	commentRepository repository_database_interfaces.CommentRepository
	userService       application_interfaces.UserService
	db                *gorm.DB
}

func NewPodcastService(
	constants *bootstrap.Constants,
	awsS3Service application_interfaces.S3Service,
	categoryService application_interfaces.CategoryService,
	podcastRepository repository_database_interfaces.PodcastRepository,
	commentRepository repository_database_interfaces.CommentRepository,
	userService application_interfaces.UserService,
	db *gorm.DB,

) *PodcastService {
	return &PodcastService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		categoryService:   categoryService,
		podcastRepository: podcastRepository,
		commentRepository: commentRepository,
		userService:       userService,
		db:                db,
	}
}

func (podcastService *PodcastService) fetchPodcastByID(podcastID uint) *entities.Podcast {
	var notFoundError exceptions.NotFoundError
	podcast, podcastExist := podcastService.podcastRepository.FindPodcastByID(podcastService.db, podcastID)
	if !podcastExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Podcast
		panic(notFoundError)
	}
	return podcast
}

func (podcastService *PodcastService) validateUniquePodcastName(name string) {
	var conflictError exceptions.ConflictError
	_, podcastExist := podcastService.podcastRepository.FindPodcastByName(podcastService.db, name)
	if podcastExist {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Tittle,
			podcastService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (podcastService *PodcastService) setPodcastBannerPath(podcast *entities.Podcast, banner *multipart.FileHeader) {
	bannerPath := podcastService.constants.S3Service.GetPodcastBannerKey(podcast.ID, banner.Filename)
	podcastService.awsS3Service.UploadObject(enums.PodcastsBucket, bannerPath, banner)
	podcast.BannerPath = bannerPath
}

func (podcastService *PodcastService) GetPodcastList(page, pageSize int) []dto.PodcastDetailsResponse {
	offset := (page - 1) * pageSize
	podcasts, _ := podcastService.podcastRepository.FindAllPodcasts(podcastService.db, offset, pageSize)
	podcastsDetails := make([]dto.PodcastDetailsResponse, len(podcasts))
	for i, podcast := range podcasts {
		banner := ""
		if podcast.BannerPath != "" {
			banner = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, podcast.BannerPath, 8*time.Hour)
		}
		publisher, _ := podcastService.userService.FindByUserID(podcast.PublisherID)
		podcastsDetails[i] = dto.PodcastDetailsResponse{
			ID:               podcast.ID,
			CreatedAt:        podcast.CreatedAt,
			Name:             podcast.Name,
			Description:      podcast.Description,
			Banner:           banner,
			Publisher:        publisher.Name,
			SubscribersCount: len(podcast.Subscribers),
		}
	}
	return podcastsDetails
}

func (podcastService *PodcastService) GetPodcastDetails(podcastID uint) dto.PodcastDetailsResponse {
	var notFoundError exceptions.NotFoundError

	podcast, podcastExist := podcastService.podcastRepository.FindDetailedPodcastByID(podcastService.db, podcastID)
	if !podcastExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Podcast
		panic(notFoundError)
	}

	banner := ""
	if podcast.BannerPath != "" {
		banner = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, podcast.BannerPath, 8*time.Hour)
	}

	publisher, _ := podcastService.userService.FindByUserID(podcast.PublisherID)

	categories := make([]string, len(podcast.Categories))
	for i, category := range podcast.Categories {
		categories[i] = category.Name
	}

	podcastDetails := dto.PodcastDetailsResponse{
		ID:               podcastID,
		CreatedAt:        podcast.CreatedAt,
		Name:             podcast.Name,
		Description:      podcast.Description,
		Banner:           banner,
		Publisher:        publisher.Name,
		Categories:       categories,
		SubscribersCount: len(podcast.Subscribers),
	}

	return podcastDetails
}

func (podcastService *PodcastService) CreatePodcast(
	name, description string,
	categoryNames []string,
	banner *multipart.FileHeader,
	publisherID uint,
) *entities.Podcast {
	podcastService.validateUniquePodcastName(name)
	categories := podcastService.categoryService.GetCategoriesByName(categoryNames)

	var podcast *entities.Podcast
	err := repository_database.ExecuteInTransaction(podcastService.db, func(tx *gorm.DB) error {
		commentable := podcastService.commentRepository.CreateNewCommentable(tx)

		podcast = &entities.Podcast{
			ID:          commentable.CID,
			Name:        name,
			Description: description,
			PublisherID: publisherID,
			Categories:  categories,
		}
		if err := podcastService.podcastRepository.CreatePodcast(tx, podcast); err != nil {
			panic(err)
		}

		if banner != nil {
			podcastService.setPodcastBannerPath(podcast, banner)
			if err := podcastService.podcastRepository.UpdatePodcast(tx, podcast); err != nil {
				panic(err)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return podcast
}

func (podcastService *PodcastService) UpdatePodcast(podcastID uint, name, description *string, categories *[]string, banner *multipart.FileHeader) {
	podcast := podcastService.fetchPodcastByID(podcastID)

	err := repository_database.ExecuteInTransaction(podcastService.db, func(tx *gorm.DB) error {
		if name != nil {
			podcastService.validateUniquePodcastName(*name)
			podcast.Name = *name
		}
		if description != nil {
			podcast.Description = *description
		}
		if categories != nil {
			categoryModels := podcastService.categoryService.GetCategoriesByName(*categories)
			podcastService.podcastRepository.UpdatePodcastCategories(tx, podcastID, categoryModels)
		}
		if banner != nil {
			if podcast.BannerPath != "" {
				podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, podcast.BannerPath)
			}
			podcastService.setPodcastBannerPath(podcast, banner)
		}
		if err := podcastService.podcastRepository.UpdatePodcast(tx, podcast); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (podcastService *PodcastService) DeletePodcast(podcastID uint) {
	podcast := podcastService.fetchPodcastByID(podcastID)

	err := repository_database.ExecuteInTransaction(podcastService.db, func(tx *gorm.DB) error {
		if err := podcastService.podcastRepository.DeletePodcast(tx, podcast.ID); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, podcast.BannerPath)
	for _, podcast := range podcast.Episodes {
		podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, podcast.AudioPath)
	}
}

func (podcastService *PodcastService) SubscribePodcast(podcastID, userID uint) {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError

	podcast := podcastService.fetchPodcastByID(podcastID)

	user, userExist := podcastService.userService.FindByUserID(userID)
	if !userExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.User
		panic(notFoundError)
	}

	if podcastService.podcastRepository.ExistSubscriberByID(podcastService.db, podcast, userID) {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Podcast,
			podcastService.constants.ErrorTag.AlreadySubscribed)
		panic(conflictError)
	}

	podcastService.podcastRepository.SubscribePodcast(podcastService.db, podcast, user)
}

func (podcastService *PodcastService) UnSubscribePodcast(podcastID, userID uint) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError

	podcast := podcastService.fetchPodcastByID(podcastID)

	user, userExist := podcastService.userService.FindByUserID(userID)
	if !userExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.User
		panic(notFoundError)
	}

	if !podcastService.podcastRepository.ExistSubscriberByID(podcastService.db, podcast, userID) {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Podcast,
			podcastService.constants.ErrorTag.NotSubscribe)
		panic(conflictError)
	}

	podcastService.podcastRepository.UnSubscribePodcast(podcastService.db, podcast, user)
}

func (podcastService *PodcastService) IsUserSubscribedPodcast(podcastID, userID uint) bool {
	var notFoundError exceptions.NotFoundError

	podcast := podcastService.fetchPodcastByID(podcastID)

	_, userExist := podcastService.userService.FindByUserID(userID)
	if !userExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.User
		panic(notFoundError)
	}

	return podcastService.podcastRepository.ExistSubscriberByID(podcastService.db, podcast, userID)
}

func (podcastService *PodcastService) GetEpisodesList(page, pageSize int) []dto.EpisodeDetailsResponse {
	offset := (page - 1) * pageSize
	episodes, _ := podcastService.podcastRepository.FindAllEpisodes(podcastService.db, offset, pageSize)
	episodesDetails := make([]dto.EpisodeDetailsResponse, len(episodes))
	for i, episode := range episodes {
		banner := ""
		if episode.BannerPath != "" {
			banner = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, episode.BannerPath, 8*time.Hour)
		}
		audio := ""
		if episode.AudioPath != "" {
			audio = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, episode.AudioPath, 8*time.Hour)
		}
		publisher, _ := podcastService.userService.FindByUserID(episode.PublisherID)
		episodesDetails[i] = dto.EpisodeDetailsResponse{
			ID:          episode.ID,
			CreatedAt:   episode.CreatedAt,
			Name:        episode.Name,
			Description: episode.Description,
			Banner:      banner,
			Audio:       audio,
			Publisher:   publisher.Name,
		}
	}
	return episodesDetails
}

func (podcastService *PodcastService) GetEpisodeDetails(episodeID uint) dto.EpisodeDetailsResponse {
	var notFoundError exceptions.NotFoundError
	episode, episodeExist := podcastService.podcastRepository.FindEpisodeByID(podcastService.db, episodeID)
	if !episodeExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Episode
		panic(notFoundError)
	}

	banner := ""
	if episode.BannerPath != "" {
		banner = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, episode.BannerPath, 8*time.Hour)
	}
	audio := ""
	if episode.AudioPath != "" {
		audio = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, episode.AudioPath, 8*time.Hour)
	}
	publisher, _ := podcastService.userService.FindByUserID(episode.PublisherID)
	episodeDetails := dto.EpisodeDetailsResponse{
		ID:          episode.ID,
		CreatedAt:   episode.CreatedAt,
		Name:        episode.Name,
		Description: episode.Description,
		Banner:      banner,
		Audio:       audio,
		Publisher:   publisher.Name,
	}
	return episodeDetails
}

func (podcastService *PodcastService) fetchEpisodeByID(episodeID uint) *entities.Episode {
	var notFoundError exceptions.NotFoundError
	episode, episodeExist := podcastService.podcastRepository.FindEpisodeByID(podcastService.db, episodeID)
	if !episodeExist {
		notFoundError.ErrorField = podcastService.constants.ErrorField.Episode
		panic(notFoundError)
	}
	return episode
}

func (podcastService *PodcastService) validateUniqueEpisodeName(name string, podcastID uint) {
	var conflictError exceptions.ConflictError
	_, episodeExist := podcastService.podcastRepository.FindPodcastEpisodeByName(podcastService.db, name, podcastID)
	if episodeExist {
		conflictError.AppendError(
			podcastService.constants.ErrorField.Tittle,
			podcastService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (podcastService *PodcastService) setEpisodeBannerPath(podcastID uint, episode *entities.Episode, banner *multipart.FileHeader) {
	bannerPath := podcastService.constants.S3Service.GetPodcastEpisodeBannerKey(podcastID, episode.ID, banner.Filename)
	podcastService.awsS3Service.UploadObject(enums.PodcastsBucket, bannerPath, banner)
	episode.BannerPath = bannerPath
}

func (podcastService *PodcastService) setEpisodeAudioPath(podcastID uint, episode *entities.Episode, audio *multipart.FileHeader) {
	audioPath := podcastService.constants.S3Service.GetPodcastEpisodeKey(podcastID, episode.ID, audio.Filename)
	podcastService.awsS3Service.UploadObject(enums.PodcastsBucket, audioPath, audio)
	episode.AudioPath = audioPath
}

func (podcastService *PodcastService) CreateEpisode(name, description string, banner, audio *multipart.FileHeader, podcastID, publisherID uint) *entities.Episode {
	podcastService.fetchPodcastByID(podcastID)
	podcastService.validateUniqueEpisodeName(name, podcastID)

	var episode *entities.Episode
	err := repository_database.ExecuteInTransaction(podcastService.db, func(tx *gorm.DB) error {
		episode = &entities.Episode{
			Name:        name,
			Description: description,
			PublisherID: publisherID,
			PodcastID:   podcastID,
		}
		if err := podcastService.podcastRepository.CreateEpisode(tx, episode); err != nil {
			panic(err)
		}

		if banner != nil {
			podcastService.setEpisodeBannerPath(podcastID, episode, banner)
		}

		podcastService.setEpisodeAudioPath(podcastID, episode, audio)

		if err := podcastService.podcastRepository.UpdateEpisode(tx, episode); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return episode
}

func (podcastService *PodcastService) UpdateEpisode(episodeID uint, name, description *string, banner, audio *multipart.FileHeader) {
	episode := podcastService.fetchEpisodeByID(episodeID)

	err := repository_database.ExecuteInTransaction(podcastService.db, func(tx *gorm.DB) error {
		if name != nil {
			podcastService.validateUniqueEpisodeName(*name, episode.PodcastID)
			episode.Name = *name
		}

		if description != nil {
			episode.Description = *description
		}

		if banner != nil {
			if episode.BannerPath != "" {
				podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, episode.BannerPath)
			}
			podcastService.setEpisodeBannerPath(episode.PodcastID, episode, banner)
		}

		if audio != nil {
			podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, episode.AudioPath)
			podcastService.setEpisodeAudioPath(episode.PodcastID, episode, audio)
		}

		if err := podcastService.podcastRepository.UpdateEpisode(tx, episode); err != nil {
			panic(err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (podcastService *PodcastService) DeleteEpisode(episodeID uint) {
	episode := podcastService.fetchEpisodeByID(episodeID)

	err := repository_database.ExecuteInTransaction(podcastService.db, func(tx *gorm.DB) error {
		if err := podcastService.podcastRepository.DeleteEpisode(tx, episodeID); err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	if episode.BannerPath != "" {
		podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, episode.BannerPath)
	}
	if episode.AudioPath != "" {
		podcastService.awsS3Service.DeleteObject(enums.PodcastsBucket, episode.AudioPath)
	}
}

func (podcastService *PodcastService) SearchEvents(query string, page, pageSize int) []dto.PodcastDetailsResponse {
	var podcasts []*entities.Podcast
	offset := (page - 1) * pageSize
	if query != "" {
		podcasts = podcastService.podcastRepository.FullTextSearch(podcastService.db, query, offset, pageSize)
	} else {
		podcasts, _ = podcastService.podcastRepository.FindAllPodcasts(podcastService.db, offset, pageSize)
	}

	podcastsDetails := make([]dto.PodcastDetailsResponse, len(podcasts))
	for i, podcast := range podcasts {
		banner := ""
		if podcast.BannerPath != "" {
			banner = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, podcast.BannerPath, 8*time.Hour)
		}
		publisher, _ := podcastService.userService.FindByUserID(podcast.PublisherID)
		podcastsDetails[i] = dto.PodcastDetailsResponse{
			ID:               podcast.ID,
			CreatedAt:        podcast.CreatedAt,
			Name:             podcast.Name,
			Description:      podcast.Description,
			Banner:           banner,
			Publisher:        publisher.Name,
			SubscribersCount: len(podcast.Subscribers),
		}
	}
	return podcastsDetails
}

func (podcastService *PodcastService) FilterPodcastsByCategory(categories []string, page, pageSize int) []dto.PodcastDetailsResponse {
	var podcasts []*entities.Podcast
	offset := (page - 1) * pageSize
	if len(categories) == 0 {
		podcasts, _ = podcastService.podcastRepository.FindAllPodcasts(podcastService.db, offset, pageSize)
	} else {
		podcasts = podcastService.podcastRepository.FindPodcastsByCategoryName(podcastService.db, categories, offset, pageSize)
	}

	podcastsDetails := make([]dto.PodcastDetailsResponse, len(podcasts))
	for i, podcast := range podcasts {
		banner := ""
		if podcast.BannerPath != "" {
			banner = podcastService.awsS3Service.GetPresignedURL(enums.PodcastsBucket, podcast.BannerPath, 8*time.Hour)
		}
		publisher, _ := podcastService.userService.FindByUserID(podcast.PublisherID)
		podcastsDetails[i] = dto.PodcastDetailsResponse{
			ID:               podcast.ID,
			CreatedAt:        podcast.CreatedAt,
			Name:             podcast.Name,
			Description:      podcast.Description,
			Banner:           banner,
			Publisher:        publisher.Name,
			SubscribersCount: len(podcast.Subscribers),
		}
	}
	return podcastsDetails
}
