package repository_database

import (
	"first-project/src/entities"

	"gorm.io/gorm"
)

type PodcastRepository struct {
	db *gorm.DB
}

func NewPodcastRepository(db *gorm.DB) *PodcastRepository {
	return &PodcastRepository{
		db: db,
	}
}

func (repo *PodcastRepository) FindCategoriesByNames(categoryNames []string) []entities.Category {
	var categories []entities.Category

	for _, categoryName := range categoryNames {
		var category entities.Category
		if err := repo.db.FirstOrCreate(&category, entities.Category{Name: categoryName}).Error; err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}
	return categories
}

func (repo *PodcastRepository) FindPodcastByName(name string) (*entities.Podcast, bool) {
	var podcast entities.Podcast
	result := repo.db.First(&podcast, "name = ?", name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &podcast, false
		}
		panic(result.Error)
	}
	return &podcast, true
}

func (repo *PodcastRepository) FindPodcastByID(podcastID uint) (*entities.Podcast, bool) {
	var podcast entities.Podcast
	result := repo.db.First(&podcast, "id = ?", podcastID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &podcast, false
		}
		panic(result.Error)
	}
	return &podcast, true
}

func (repo *PodcastRepository) CreatePodcast(podcast *entities.Podcast) *entities.Podcast {
	result := repo.db.Create(podcast)
	if result.Error != nil {
		panic(result.Error)
	}
	return podcast
}

func (repo *PodcastRepository) UpdatePodcast(podcast *entities.Podcast) {
	err := repo.db.Save(podcast).Error
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) FindEpisodeByID(episodeID uint) (*entities.Episode, bool) {
	var episode entities.Episode
	result := repo.db.First(&episode, "id = ?", episodeID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &episode, false
		}
		panic(result.Error)
	}
	return &episode, true
}

func (repo *PodcastRepository) FindEpisodeByName(name string) (*entities.Episode, bool) {
	var episode entities.Episode
	result := repo.db.First(&episode, "name = ?", name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &episode, false
		}
		panic(result.Error)
	}
	return &episode, true
}

func (repo *PodcastRepository) CreateEpisode(episode *entities.Episode) *entities.Episode {
	result := repo.db.Create(episode)
	if result.Error != nil {
		panic(result.Error)
	}
	return episode
}

func (repo *PodcastRepository) UpdateEpisode(episode *entities.Episode) {
	err := repo.db.Save(episode).Error
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) DeleteEpisodeByID(episodeID uint) {
	err := repo.db.Unscoped().Delete(&entities.Episode{}, episodeID).Error
	if err != nil {
		panic(err)
	}
}
