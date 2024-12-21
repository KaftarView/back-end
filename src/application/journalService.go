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
		// BannerPath: ,
		// JournalFilePath: ,
		AuthorID: authorID,
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
