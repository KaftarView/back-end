package application

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"mime/multipart"
	"time"
)

type journalService struct {
	constants         *bootstrap.Constants
	awsS3Service      *application_aws.S3service
	userRepository    repository_database_interfaces.UserRepository
	journalRepository repository_database_interfaces.JournalRepository
}

func NewJournalService(
	constants *bootstrap.Constants,
	awsS3Service *application_aws.S3service,
	userRepository repository_database_interfaces.UserRepository,
	journalRepository repository_database_interfaces.JournalRepository,
) *journalService {
	return &journalService{
		constants:         constants,
		awsS3Service:      awsS3Service,
		userRepository:    userRepository,
		journalRepository: journalRepository,
	}
}

func (journalService *journalService) GetJournalsList(page, pageSize int) []dto.JournalDetailsResponse {
	offset := (page - 1) * pageSize
	journalsList, _ := journalService.journalRepository.FindAllJournals(offset, pageSize)

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

func (journalService *journalService) CreateJournal(name, description string, banner, journalFile *multipart.FileHeader, authorID uint) *entities.Journal {
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

	bannerPath := journalService.constants.S3Service.GetJournalBannerKey(journalModel.ID, banner.Filename)
	journalService.awsS3Service.UploadObject(enums.JournalsBucket, bannerPath, banner)
	filePath := journalService.constants.S3Service.GetJournalFileKey(journalModel.ID, journalFile.Filename)
	journalService.awsS3Service.UploadObject(enums.JournalsBucket, filePath, journalFile)

	journalModel.BannerPath = bannerPath
	journalModel.JournalFilePath = filePath
	journalService.journalRepository.UpdateJournal(journalModel)

	return journalModel
}

func (journalService *journalService) UpdateJournal(journalID uint, name, description *string, banner, journalFile *multipart.FileHeader) {
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
			journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.BannerPath)
		}
		bannerPath := journalService.constants.S3Service.GetJournalBannerKey(journalID, banner.Filename)
		journalService.awsS3Service.UploadObject(enums.JournalsBucket, bannerPath, banner)
		journal.BannerPath = bannerPath
	}
	if journalFile != nil {
		if journal.JournalFilePath != "" {
			journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.JournalFilePath)
		}
		filePath := journalService.constants.S3Service.GetJournalFileKey(journalID, journalFile.Filename)
		journalService.awsS3Service.UploadObject(enums.JournalsBucket, filePath, journalFile)
		journal.JournalFilePath = filePath
	}

	journalService.journalRepository.UpdateJournal(journal)
}

func (journalService *journalService) DeleteJournal(journalID uint) {
	var notFoundError exceptions.NotFoundError
	journal, journalExist := journalService.journalRepository.FindJournalByID(journalID)
	if !journalExist {
		notFoundError.ErrorField = journalService.constants.ErrorField.Journal
		panic(notFoundError)
	}

	if journal.BannerPath != "" {
		journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.BannerPath)
	}
	if journal.JournalFilePath != "" {
		journalService.awsS3Service.DeleteObject(enums.JournalsBucket, journal.JournalFilePath)
	}

	journalService.journalRepository.DeleteJournal(journalID)
}

func (journalService *journalService) SearchJournals(query string, page, pageSize int) []dto.JournalDetailsResponse {
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
			banner = journalService.awsS3Service.GetPresignedURL(enums.JournalsBucket, journal.BannerPath, 8*time.Hour)
		}
		file := ""
		if journal.JournalFilePath != "" {
			banner = journalService.awsS3Service.GetPresignedURL(enums.JournalsBucket, journal.JournalFilePath, 8*time.Hour)
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
