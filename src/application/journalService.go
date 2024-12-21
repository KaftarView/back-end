package application

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	repository_database "first-project/src/repository/database"
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
