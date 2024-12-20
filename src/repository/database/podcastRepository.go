package repository_database

import (
	"first-project/src/entities"
	"strings"

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

func (repo *PodcastRepository) BeginTransaction() *gorm.DB {
	tx := repo.db.Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}
	return tx
}

func (repo *PodcastRepository) CreatePodcast(tx *gorm.DB, podcast *entities.Podcast) *entities.Podcast {
	result := tx.Create(podcast)
	if result.Error != nil {
		panic(result.Error)
	}
	return podcast
}

func (repo *PodcastRepository) UpdatePodcast(tx *gorm.DB, podcast *entities.Podcast) {
	err := tx.Save(podcast).Error
	if err != nil {
		panic(err)
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

func (repo *PodcastRepository) FindAllPodcasts() ([]*entities.Podcast, bool) {
	var podcasts []*entities.Podcast
	result := repo.db.
		Preload("Subscribers").
		Find(&podcasts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return podcasts, false
		}
		panic(result.Error)
	}
	return podcasts, true
}

func (repo *PodcastRepository) FindDetailedPodcastByID(podcastID uint) (*entities.Podcast, bool) {
	var podcast entities.Podcast
	result := repo.db.
		Preload("Subscribers").
		Preload("Episodes").
		Preload("Categories").
		First(&podcast, podcastID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}

	return &podcast, true
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

func (repo *PodcastRepository) ExistSubscriberByID(podcast *entities.Podcast, userID uint) bool {
	var count int64
	result := repo.db.Model(podcast).
		Joins("JOIN podcast_subscribers ON podcast_subscribers.user_id = ? AND podcast_subscribers.podcast_id = ?", userID, podcast.ID).
		Count(&count)
	if result.Error != nil {
		panic(result.Error)
	}
	return count > 0
}

func (repo *PodcastRepository) DeletePodcast(podcastID uint) {
	err := repo.db.Unscoped().Delete(&entities.Podcast{}, podcastID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) SubscribePodcast(podcast *entities.Podcast, user *entities.User) {
	err := repo.db.Model(podcast).Association("Subscribers").Append(user)
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) UnSubscribePodcast(podcast *entities.Podcast, user *entities.User) {
	err := repo.db.Model(&podcast).Association("Subscribers").Delete(user)
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

func (repo *PodcastRepository) FindPodcastEpisodeByName(name string, podcastID uint) (*entities.Episode, bool) {
	var episode entities.Episode
	result := repo.db.First(&episode, "name = ? AND podcast_id = ?", name, podcastID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &episode, false
		}
		panic(result.Error)
	}
	return &episode, true
}

func (repo *PodcastRepository) FindAllEpisodes() ([]*entities.Episode, bool) {
	var podcasts []*entities.Episode
	result := repo.db.Find(&podcasts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return podcasts, false
		}
		panic(result.Error)
	}
	return podcasts, true
}

func (repo *PodcastRepository) CreateEpisode(tx *gorm.DB, episode *entities.Episode) *entities.Episode {
	result := tx.Create(episode)
	if result.Error != nil {
		panic(result.Error)
	}
	return episode
}

func (repo *PodcastRepository) UpdateEpisode(tx *gorm.DB, episode *entities.Episode) {
	err := tx.Save(episode).Error
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) DeleteEpisodeByID(tx *gorm.DB, episodeID uint) {
	err := tx.Unscoped().Delete(&entities.Episode{}, episodeID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) FullTextSearch(query string, offset, pageSize int) []*entities.Podcast {
	var podcasts []*entities.Podcast

	repo.db.Exec(`ALTER TABLE podcasts ADD FULLTEXT INDEX idx_name_description (name, description)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := repo.db.Model(&entities.Podcast{}).
		Where("MATCH(name, description) AGAINST(? IN BOOLEAN MODE)", searchQuery).
		Offset(offset).
		Limit(pageSize).
		Find(&podcasts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return podcasts
}

func (repo *PodcastRepository) FindPodcastsByCategoryName(categories []string, offset, pageSize int) []*entities.Podcast {
	var podcasts []*entities.Podcast

	result := repo.db.
		Distinct("podcasts.*").
		Joins("JOIN podcast_categories ON podcasts.id = podcast_categories.podcast_id").
		Joins("JOIN categories ON categories.id = podcast_categories.category_id").
		Where("categories.name IN ?", categories).
		Limit(pageSize).
		Offset(offset).
		Find(&podcasts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}

	return podcasts
}
