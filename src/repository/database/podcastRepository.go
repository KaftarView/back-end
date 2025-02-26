package repository_database

import (
	"first-project/src/entities"
	"strings"

	"gorm.io/gorm"
)

type PodcastRepository struct{}

func NewPodcastRepository() *PodcastRepository {
	return &PodcastRepository{}
}

func (repo *PodcastRepository) CreatePodcast(db *gorm.DB, podcast *entities.Podcast) error {
	result := db.Create(podcast)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *PodcastRepository) UpdatePodcastCategories(db *gorm.DB, podcastID uint, categories []entities.Category) {
	err := db.Model(&entities.Podcast{ID: podcastID}).Association("Categories").Replace(categories)
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) UpdatePodcast(db *gorm.DB, podcast *entities.Podcast) error {
	err := db.Save(podcast).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *PodcastRepository) FindAllPodcasts(db *gorm.DB, offset, pageSize int) ([]*entities.Podcast, bool) {
	var podcasts []*entities.Podcast
	result := OrderByCreatedAtDesc(db).
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

func (repo *PodcastRepository) FindDetailedPodcastByID(db *gorm.DB, podcastID uint) (*entities.Podcast, bool) {
	var podcast entities.Podcast
	result := db.
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

func (repo *PodcastRepository) FindPodcastByName(db *gorm.DB, name string) (*entities.Podcast, bool) {
	var podcast entities.Podcast
	result := db.First(&podcast, "name = ?", name)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &podcast, false
		}
		panic(result.Error)
	}
	return &podcast, true
}

func (repo *PodcastRepository) FindPodcastByID(db *gorm.DB, podcastID uint) (*entities.Podcast, bool) {
	var podcast entities.Podcast
	result := db.First(&podcast, "id = ?", podcastID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &podcast, false
		}
		panic(result.Error)
	}
	return &podcast, true
}

func (repo *PodcastRepository) ExistSubscriberByID(db *gorm.DB, podcast *entities.Podcast, userID uint) bool {
	var count int64
	result := db.Model(podcast).
		Joins("JOIN podcast_subscribers ON podcast_subscribers.user_id = ? AND podcast_subscribers.podcast_id = ?", userID, podcast.ID).
		Count(&count)
	if result.Error != nil {
		panic(result.Error)
	}
	return count > 0
}

func (repo *PodcastRepository) DeletePodcast(db *gorm.DB, podcastID uint) error {
	err := db.Unscoped().Delete(&entities.Podcast{}, podcastID).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *PodcastRepository) SubscribePodcast(db *gorm.DB, podcast *entities.Podcast, user *entities.User) {
	err := db.Model(podcast).Association("Subscribers").Append(user)
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) UnSubscribePodcast(db *gorm.DB, podcast *entities.Podcast, user *entities.User) {
	err := db.Model(&podcast).Association("Subscribers").Delete(user)
	if err != nil {
		panic(err)
	}
}

func (repo *PodcastRepository) FindEpisodeByID(db *gorm.DB, episodeID uint) (*entities.Episode, bool) {
	var episode entities.Episode
	result := db.First(&episode, "id = ?", episodeID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &episode, false
		}
		panic(result.Error)
	}
	return &episode, true
}

func (repo *PodcastRepository) FindPodcastEpisodeByName(db *gorm.DB, name string, podcastID uint) (*entities.Episode, bool) {
	var episode entities.Episode
	result := db.First(&episode, "name = ? AND podcast_id = ?", name, podcastID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &episode, false
		}
		panic(result.Error)
	}
	return &episode, true
}

func (repo *PodcastRepository) FindAllEpisodes(db *gorm.DB, offset, pageSize int) ([]*entities.Episode, bool) {
	var podcasts []*entities.Episode
	result := OrderByCreatedAtDesc(db).Offset(offset).Limit(pageSize).Find(&podcasts)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return podcasts, false
		}
		panic(result.Error)
	}
	return podcasts, true
}

func (repo *PodcastRepository) CreateEpisode(db *gorm.DB, episode *entities.Episode) error {
	return db.Create(episode).Error
}

func (repo *PodcastRepository) UpdateEpisode(db *gorm.DB, episode *entities.Episode) error {
	return db.Save(episode).Error
}

func (repo *PodcastRepository) DeleteEpisode(db *gorm.DB, episodeID uint) error {
	return db.Unscoped().Delete(&entities.Episode{}, episodeID).Error
}

func (repo *PodcastRepository) FullTextSearch(db *gorm.DB, query string, offset, pageSize int) []*entities.Podcast {
	var podcasts []*entities.Podcast

	db.Exec(`ALTER TABLE podcasts ADD FULLTEXT INDEX idx_name_description (name, description)`)
	searchQuery := "+" + strings.Join(strings.Fields(query), "* +") + "*"

	result := OrderByCreatedAtDesc(db).
		Model(&entities.Podcast{}).
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

func (repo *PodcastRepository) FindPodcastsByCategoryName(db *gorm.DB, categories []string, offset, pageSize int) []*entities.Podcast {
	var podcasts []*entities.Podcast
	result := OrderByCreatedAtDesc(db).
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
