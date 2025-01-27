package application

import (
	application_aws "first-project/src/application/aws"
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

type journalService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3service
	userService       application_interfaces.UserService
	journalRepository repository_database_interfaces.JournalRepository
	db                *gorm.DB
}

func NewJournalService(
	constants *bootstrap.Constants,
	awsS3Service *application_aws.S3service,
	userService application_interfaces.UserService,
	journalRepository repository_database_interfaces.JournalRepository,
	db *gorm.DB,
) *journalService {
	return &journalService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		userService:       userService,
		journalRepository: journalRepository,
		db:                db,
	}
}

func (journalService *journalService) fetchJournalByID(journalID uint) *entities.Journal {
	var notFoundError exceptions.NotFoundError
	journal, journalExist := journalService.journalRepository.FindJournalByID(journalService.db, journalID)
	if !journalExist {
		notFoundError.ErrorField = journalService.constants.ErrorField.Journal
		panic(notFoundError)
	}
	return journal
}

func (journalService *journalService) validateUniqueJournalName(name string) {
	var conflictError exceptions.ConflictError
	_, journalExist := journalService.journalRepository.FindJournalByName(journalService.db, name)
	if journalExist {
		conflictError.AppendError(
			journalService.constants.ErrorField.Journal,
			journalService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}
}

func (journalService *journalService) setJournalBannerPath(journal *entities.Journal, banner *multipart.FileHeader) {
	bannerPath := journalService.constants.S3Service.GetJournalBannerKey(journal.ID, banner.Filename)
	journalService.awsS3Service.UploadObject(enums.JournalsBucket, bannerPath, banner)
	journal.BannerPath = bannerPath
}

func (journalService *journalService) setJournalFilePath(journal *entities.Journal, file *multipart.FileHeader) {
	filePath := journalService.constants.S3Service.GetJournalFileKey(journal.ID, file.Filename)
	journalService.awsS3Service.UploadObject(enums.JournalsBucket, filePath, file)
	journal.JournalFilePath = filePath
}

func (journalService *journalService) GetJournalsList(page, pageSize int) []dto.JournalDetailsResponse {
	offset := (page - 1) * pageSize
	journalsList, _ := journalService.journalRepository.FindAllJournals(journalService.db, offset, pageSize)

	journalsDetails := make([]dto.JournalDetailsResponse, len(journalsList))
	for i, journal := range journalsList {
		banner := ""
		if journal.BannerPath != "" {
			banner = journalService.awsS3Service.GetPresignedURL(enums.JournalsBucket, journal.BannerPath, 8*time.Hour)
		}
		file := ""
		if journal.JournalFilePath != "" {
			file = journalService.awsS3Service.GetPresignedURL(enums.JournalsBucket, journal.JournalFilePath, 8*time.Hour)
		}
		author, _ := journalService.userService.FindByUserID(journal.AuthorID)
		journalsDetails[i] = dto.JournalDetailsResponse{
			ID:          journal.ID,
			Name:        journal.Name,
			CreatedAt:   journal.CreatedAt,
			Description: journal.Description,
			Banner:      banner,
			JournalFile: file,
			Author:      author.Name,
		}
	}

	return journalsDetails
}

func (journalService *journalService) CreateJournal(name, description string, banner, journalFile *multipart.FileHeader, authorID uint) *entities.Journal {
	journalService.validateUniqueJournalName(name)

	journal := &entities.Journal{
		Name:        name,
		Description: description,
		AuthorID:    authorID,
	}
	err := repository_database.ExecuteInTransaction(journalService.db, func(tx *gorm.DB) error {
		if err := journalService.journalRepository.CreateJournal(tx, journal); err != nil {
			panic(err)
		}

		journalService.setJournalBannerPath(journal, banner)
		journalService.setJournalFilePath(journal, journalFile)

		if err := journalService.journalRepository.UpdateJournal(tx, journal); err != nil {
			panic(err)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return journal
}

func (journalService *journalService) UpdateJournal(journalID uint, name, description *string, banner, journalFile *multipart.FileHeader) {
	journal := journalService.fetchJournalByID(journalID)

	err := repository_database.ExecuteInTransaction(journalService.db, func(tx *gorm.DB) error {
		if name != nil {
			journalService.validateUniqueJournalName(*name)
			journal.Name = *name
		}
		if description != nil {
			journal.Description = *description
		}
		if banner != nil {
			if journal.BannerPath != "" {
				journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.BannerPath)
			}
			journalService.setJournalBannerPath(journal, banner)
		}
		if journalFile != nil {
			if journal.JournalFilePath != "" {
				journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.JournalFilePath)
			}
			journalService.setJournalFilePath(journal, journalFile)
		}

		journalService.journalRepository.UpdateJournal(tx, journal)

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (journalService *journalService) DeleteJournal(journalID uint) {
	journal := journalService.fetchJournalByID(journalID)

	if err := journalService.journalRepository.DeleteJournal(journalService.db, journalID); err != nil {
		panic(err)
	}
	if journal.BannerPath != "" {
		journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.BannerPath)
	}
	if journal.JournalFilePath != "" {
		journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.JournalFilePath)
	}
}

func (journalService *journalService) SearchJournals(query string, page, pageSize int) []dto.JournalDetailsResponse {
	var journalsList []*entities.Journal
	offset := (page - 1) * pageSize
	if query != "" {
		journalsList = journalService.journalRepository.FullTextSearch(journalService.db, query, offset, pageSize)
	} else {
		journalsList, _ = journalService.journalRepository.FindAllJournals(journalService.db, offset, pageSize)
	}

	journalsDetails := make([]dto.JournalDetailsResponse, len(journalsList))
	for i, journal := range journalsList {
		banner := ""
		if journal.BannerPath != "" {
			banner = journalService.awsS3Service.GetPresignedURL(enums.JournalsBucket, journal.BannerPath, 8*time.Hour)
		}
		file := ""
		if journal.JournalFilePath != "" {
			banner = journalService.awsS3Service.GetPresignedURL(enums.JournalsBucket, journal.JournalFilePath, 8*time.Hour)
		}
		author, _ := journalService.userService.FindByUserID(journal.AuthorID)
		journalsDetails[i] = dto.JournalDetailsResponse{
			ID:          journal.ID,
			Name:        journal.Name,
			CreatedAt:   journal.CreatedAt,
			Description: journal.Description,
			Banner:      banner,
			JournalFile: file,
			Author:      author.Name,
		}
	}

	return journalsDetails
}
