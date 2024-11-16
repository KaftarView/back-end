package application_cron

import (
	application_communication "first-project/src/application/communication/emailService"
	repository_database "first-project/src/repository/database"
	"time"

	"github.com/go-co-op/gocron"
)

type CronJob struct {
	userRepository *repository_database.UserRepository
	emailService   *application_communication.EmailService
}

func NewCronJob(
	userRepository *repository_database.UserRepository,
	emailService *application_communication.EmailService,
) *CronJob {
	return &CronJob{
		userRepository: userRepository,
		emailService:   emailService,
	}
}

func (cronJob *CronJob) RunCronJob() {
	jobScheduler := gocron.NewScheduler(time.UTC)
	jobScheduler.Every(1).Day().At("15:00").Do(cronJob.reminderEmailToUnverifiedUsers)
	jobScheduler.StartAsync()
}
