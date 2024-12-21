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

type JournalService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3service
	journalRepository *repository_database.JournalRepository
}

func NewJournalService(
	constants *bootstrap.Constants,
	awsS3Service *application_aws.S3service,
	journalRepository *repository_database.JournalRepository,
) *JournalService {
	return &JournalService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		journalRepository: journalRepository,
	}
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
