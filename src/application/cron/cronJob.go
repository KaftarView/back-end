package application_cron

import (
	application_communication "first-project/src/application/communication/emailService"
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

type CronJob struct {
	constants          *bootstrap.Constants
	userRepository     repository_database_interfaces.UserRepository
	purchaseRepository repository_database_interfaces.PurchaseRepository
	eventRepository    repository_database_interfaces.EventRepository
	emailService       application_interfaces.EmailService
	db                 *gorm.DB
}

func NewCronJob(
	constants *bootstrap.Constants,
	userRepository repository_database_interfaces.UserRepository,
	purchaseRepository repository_database_interfaces.PurchaseRepository,
	eventRepository repository_database_interfaces.EventRepository,
	emailService *application_communication.EmailService,
	db *gorm.DB,
) *CronJob {
	return &CronJob{
		constants:          constants,
		userRepository:     userRepository,
		purchaseRepository: purchaseRepository,
		eventRepository:    eventRepository,
		emailService:       emailService,
		db:                 db,
	}
}

func (cronJob *CronJob) RunCronJob() {
	jobScheduler := gocron.NewScheduler(time.UTC)
	jobScheduler.Every(1).Day().At("00:00").Do(cronJob.cleanupExpiredReservations)
	jobScheduler.StartAsync()
}
