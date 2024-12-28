package application

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	"fmt"
	"mime/multipart"
	"time"
)

type JournalService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3service
	userRepository    *repository_database.UserRepository
	journalRepository *repository_database.JournalRepository
}

func NewJournalService(
	constants *bootstrap.Constants,
	awsS3Service *application_aws.S3service,
	userRepository *repository_database.UserRepository,
	journalRepository *repository_database.JournalRepository,
) *JournalService {
	return &JournalService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		userRepository:    userRepository,
		journalRepository: journalRepository,
	}
}

func (journalService *JournalService) GetJournalsList(page, pageSize int) []dto.JournalDetailsResponse {
	offset := (page - 1) * pageSize
	journalsList, _ := journalService.journalRepository.FindAllJournals(offset, pageSize)

	journalsDetails := make([]dto.JournalDetailsResponse, len(journalsList))
	for i, journal := range journalsList {
		banner := ""
		if journal.BannerPath != "" {
			banner = journalService.awsS3Service.GetPresignedURL(enums.BannersBucket, journal.BannerPath, 8*time.Hour)
		}
		file := ""
		if journal.JournalFilePath != "" {
			file = journalService.awsS3Service.GetPresignedURL(enums.SessionsBucket, journal.JournalFilePath, 8*time.Hour)
		}
		author, _ := journalService.userRepository.FindByUserID(journal.AuthorID)
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

func (journalService *JournalService) CreateJournal(name, description string, banner, journalFile *multipart.FileHeader, authorID uint) *entities.Journal {
	var conflictError exceptions.ConflictError
	_, journalExist := journalService.journalRepository.FindJournalByName(name)
	if journalExist {
		conflictError.AppendError(
			journalService.constants.ErrorField.Journal,
			journalService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	journalModel := &entities.Journal{
		Name:        name,
		Description: description,
		AuthorID:    authorID,
	}
	journalService.journalRepository.CreateJournal(journalModel)

	bannerPath := fmt.Sprintf("banners/podcasts/%d/images/%s", journalModel.ID, banner.Filename)
	journalService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, banner)
	filePath := fmt.Sprintf("media/journals/%d/resources/%s", journalModel.ID, journalFile.Filename)
	journalService.awsS3Service.UploadObject(enums.SessionsBucket, filePath, journalFile)

	journalModel.BannerPath = bannerPath
	journalModel.JournalFilePath = filePath
	journalService.journalRepository.UpdateJournal(journalModel)

	return journalModel
}

func (journalService *JournalService) UpdateJournal(journalID uint, name, description *string, banner, journalFile *multipart.FileHeader) {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	journal, journalExist := journalService.journalRepository.FindJournalByID(journalID)
	if !journalExist {
		notFoundError.ErrorField = journalService.constants.ErrorField.Journal
		panic(notFoundError)
	}

	if name != nil {
		_, journalExist := journalService.journalRepository.FindJournalByName(*name)
		if journalExist {
			conflictError.AppendError(
				journalService.constants.ErrorField.Tittle,
				journalService.constants.ErrorTag.AlreadyExist)
			panic(conflictError)
		}
		journal.Name = *name
	}
	if description != nil {
		journal.Description = *description
	}
	if banner != nil {
		if journal.BannerPath != "" {
			journalService.awsS3Service.DeleteObject(enums.BannersBucket, journal.BannerPath)
		}
		bannerPath := fmt.Sprintf("banners/podcasts/%d/images/%s", journal.ID, banner.Filename)
		journalService.awsS3Service.UploadObject(enums.BannersBucket, bannerPath, banner)
		journal.BannerPath = bannerPath
	}
	if journalFile != nil {
		if journal.JournalFilePath != "" {
			journalService.awsS3Service.DeleteObject(enums.SessionsBucket, journal.JournalFilePath)
		}
		filePath := fmt.Sprintf("media/journals/%d/resources/%s", journal.ID, journalFile.Filename)
		journalService.awsS3Service.UploadObject(enums.SessionsBucket, filePath, journalFile)
		journal.JournalFilePath = filePath
	}

	journalService.journalRepository.UpdateJournal(journal)
}

func (journalService *JournalService) DeleteJournal(journalID uint) {
	var notFoundError exceptions.NotFoundError
	journal, journalExist := journalService.journalRepository.FindJournalByID(journalID)
	if !journalExist {
		notFoundError.ErrorField = journalService.constants.ErrorField.Journal
		panic(notFoundError)
	}

	if journal.BannerPath != "" {
		journalService.awsS3Service.DeleteObject(enums.BannersBucket, journal.BannerPath)
	}
	if journal.JournalFilePath != "" {
		journalService.awsS3Service.DeleteObject(enums.SessionsBucket, journal.JournalFilePath)
	}

	journalService.journalRepository.DeleteJournal(journalID)
}

func (journalService *JournalService) SearchJournals(query string, page, pageSize int) []dto.JournalDetailsResponse {
	var journalsList []*entities.Journal
	offset := (page - 1) * pageSize
	if query != "" {
		journalsList = journalService.journalRepository.FullTextSearch(query, offset, pageSize)
	} else {
		journalsList, _ = journalService.journalRepository.FindAllJournals(offset, pageSize)
	}

	journalsDetails := make([]dto.JournalDetailsResponse, len(journalsList))
	for i, journal := range journalsList {
		banner := ""
		if journal.BannerPath != "" {
			banner = journalService.awsS3Service.GetPresignedURL(enums.BannersBucket, journal.BannerPath, 8*time.Hour)
		}
		file := ""
		if journal.JournalFilePath != "" {
			banner = journalService.awsS3Service.GetPresignedURL(enums.SessionsBucket, journal.JournalFilePath, 8*time.Hour)
		}
		author, _ := journalService.userRepository.FindByUserID(journal.AuthorID)
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
