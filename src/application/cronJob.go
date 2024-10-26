package application

import (
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/entities"
	"first-project/src/repository"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type CronJob struct {
	userRepository *repository.UserRepository
	emailService   *application_communication.EmailService
}

func NewCronJob(userRepository *repository.UserRepository, emailService *application_communication.EmailService) *CronJob {
	return &CronJob{
		userRepository: userRepository,
		emailService:   emailService,
	}
}

func (cronJob *CronJob) reminderEmailToUnverifiedUsers() {
	startOfWeekAgo := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	endOfWeekAgo := startOfWeekAgo.Add(24 * time.Hour)
	users := cronJob.userRepository.FindUnverifiedUsersWeekAgo(startOfWeekAgo, endOfWeekAgo)

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user entities.User) {
			defer wg.Done()

			data := struct {
				Username string
			}{
				Username: user.Name,
			}
			cronJob.emailService.SendEmail(user.Email, "Activate your account", "activateAccount/en.html", data)
		}(user)
	}
	wg.Wait()
}

func (cronJob *CronJob) RunCronJob() {
	jobScheduler := gocron.NewScheduler(time.UTC)
	jobScheduler.Every(1).Day().At("00:00").Do(cronJob.reminderEmailToUnverifiedUsers)
	jobScheduler.StartAsync()
}
