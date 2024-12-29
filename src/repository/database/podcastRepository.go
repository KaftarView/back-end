package repository_database

import (
	"first-project/src/entities"
	"strings"

	"gorm.io/gorm"
)

type podcastRepository struct {
	db *gorm.DB
}

func NewPodcastRepository(db *gorm.DB) *podcastRepository {
	return &podcastRepository{
		db: db,
	}
}

func (repo *podcastRepository) BeginTransaction() *gorm.DB {
	tx := repo.db.Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}
	return tx
}

func (repo *podcastRepository) CreatePodcast(tx *gorm.DB, podcast *entities.Podcast) *entities.Podcast {
	result := tx.Create(podcast)
	if result.Error != nil {
		panic(result.Error)
	}
	return podcast
}

func (repo *podcastRepository) UpdatePodcastCategories(podcastID uint, categories []entities.Category) {
	err := repo.db.Model(&entities.Podcast{ID: podcastID}).Association("Categories").Replace(categories)
	if err != nil {
		panic(err)
	}
}

func (repo *podcastRepository) UpdatePodcast(tx *gorm.DB, podcast *entities.Podcast) {
	err := tx.Save(podcast).Error
	if err != nil {
		panic(err)
	}
}

func (repo *podcastRepository) FindAllPodcasts(offset, pageSize int) ([]*entities.Podcast, bool) {
	var podcasts []*entities.Podcast
	result := repo.db.
		Preload("Subscribers").
		Offset(offset).
		Limit(pageSize).
		Find(&podcasts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return podcasts, false
		}
		panic(result.Error)
	}
	return podcasts, true
}

func (repo *podcastRepository) FindDetailedPodcastByID(podcastID uint) (*entities.Podcast, bool) {
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

func (repo *podcastRepository) FindPodcastByName(name string) (*entities.Podcast, bool) {
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

func (repo *podcastRepository) FindPodcastByID(podcastID uint) (*entities.Podcast, bool) {
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

func (repo *podcastRepository) ExistSubscriberByID(podcast *entities.Podcast, userID uint) bool {
	var count int64
	result := repo.db.Model(podcast).
		Joins("JOIN podcast_subscribers ON podcast_subscribers.user_id = ? AND podcast_subscribers.podcast_id = ?", userID, podcast.ID).
		Count(&count)
	if result.Error != nil {
		panic(result.Error)
	}
	return count > 0
}

func (repo *podcastRepository) DeletePodcast(podcastID uint) {
	err := repo.db.Unscoped().Delete(&entities.Podcast{}, podcastID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *podcastRepository) SubscribePodcast(podcast *entities.Podcast, user *entities.User) {
	err := repo.db.Model(podcast).Association("Subscribers").Append(user)
	if err != nil {
		panic(err)
	}
}

func (repo *podcastRepository) UnSubscribePodcast(podcast *entities.Podcast, user *entities.User) {
	err := repo.db.Model(&podcast).Association("Subscribers").Delete(user)
	if err != nil {
		panic(err)
	}
}

func (repo *podcastRepository) FindEpisodeByID(episodeID uint) (*entities.Episode, bool) {
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

func (repo *podcastRepository) FindPodcastEpisodeByName(name string, podcastID uint) (*entities.Episode, bool) {
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

func (repo *podcastRepository) FindAllEpisodes(offset, pageSize int) ([]*entities.Episode, bool) {
	var podcasts []*entities.Episode
	result := repo.db.Offset(offset).Limit(pageSize).Find(&podcasts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return podcasts, false
		}
		panic(result.Error)
	}
	return podcasts, true
}

func (repo *podcastRepository) CreateEpisode(tx *gorm.DB, episode *entities.Episode) *entities.Episode {
	result := tx.Create(episode)
	if result.Error != nil {
		panic(result.Error)
	}
	return episode
}

func (repo *podcastRepository) UpdateEpisode(tx *gorm.DB, episode *entities.Episode) {
	err := tx.Save(episode).Error
	if err != nil {
		panic(err)
	}
}

func (repo *podcastRepository) DeleteEpisodeByID(tx *gorm.DB, episodeID uint) {
	err := tx.Unscoped().Delete(&entities.Episode{}, episodeID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *podcastRepository) FullTextSearch(query string, offset, pageSize int) []*entities.Podcast {
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

func (repo *podcastRepository) FindPodcastsByCategoryName(categories []string, offset, pageSize int) []*entities.Podcast {
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
