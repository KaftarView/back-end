package application_cron

import (
	"first-project/src/entities"
	"sync"
	"time"
)

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
