package application_cron

import (
	application_communication "first-project/src/application/communication/emailService"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
)

type CronJob struct {
	userRepository repository_database_interfaces.UserRepository
	emailService   *application_communication.EmailService
}

func NewCronJob(
	userRepository repository_database_interfaces.UserRepository,
	emailService *application_communication.EmailService,
) *CronJob {
	return &CronJob{
		userRepository: userRepository,
		emailService:   emailService,
	}
}

func (cronJob *CronJob) RunCronJob() {
	// some cron joc here
}
