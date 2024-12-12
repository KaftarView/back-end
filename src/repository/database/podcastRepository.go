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

func (repo *PodcastRepository) FindPodcastByName(name string) (entities.Podcast, bool) {
	var podcast entities.Podcast
	result := repo.db.First(&podcast, "name = ?", name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return podcast, false
		}
		panic(result.Error)
	}
	return podcast, true
}

func (repo *PodcastRepository) FindPodcastByID(podcastID uint) (entities.Podcast, bool) {
	var podcast entities.Podcast
	result := repo.db.First(&podcast, "id = ?", podcastID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return podcast, false
		}
		panic(result.Error)
	}
	return podcast, true
}

func (repo *PodcastRepository) CreatePodcast(podcast entities.Podcast) entities.Podcast {
	result := repo.db.Create(&podcast)
	if result.Error != nil {
		panic(result.Error)
	}
	return podcast
}

func (repo *PodcastRepository) SetPodcastBanner(bannerPath string, podcast entities.Podcast) {
	err := repo.db.Model(&podcast).Update("banner_path", bannerPath).Error
	if err != nil {
		panic(err)
	}
}
