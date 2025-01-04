package application_cron

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"

	"gorm.io/gorm"
)

func (cronJob *CronJob) processReservation(tx *gorm.DB, reservation *entities.Reservation) error {
	var notFoundError exceptions.NotFoundError
	for _, reservationItem := range reservation.Items {
		ticket, ticketExist := cronJob.eventRepository.FindEventTicketByID(cronJob.db, reservationItem.TicketID)
		if !ticketExist {
			continue
		}

		ticket.SoldCount--
		if err := cronJob.eventRepository.UpdateEventTicket(tx, ticket); err != nil {
			continue
		}
	}

	if reservation.DiscountID != nil {
		discount, discountExist := cronJob.eventRepository.FindDiscountByDiscountID(tx, *reservation.DiscountID)
		if !discountExist {
			notFoundError.ErrorField = cronJob.constants.ErrorField.Ticket
			return notFoundError
		}

		discount.UsedCount--
		if err := cronJob.eventRepository.UpdateEventDiscount(tx, discount); err != nil {
			notFoundError.ErrorField = cronJob.constants.ErrorField.Discount
			return notFoundError
		}
	}

	reservation.Status = enums.Expired
	cronJob.reservationRepository.UpdateReservation(tx, reservation)
	return nil
}

func (cronJob *CronJob) cleanupExpiredReservations() {
	err := repository_database.ExecuteInTransaction(cronJob.db, func(tx *gorm.DB) error {
		expiredReservations, err := cronJob.reservationRepository.GetExpiredReservations(tx)
		if err != nil {
			panic(err)
		}

		for _, reservation := range expiredReservations {
			if err := cronJob.processReservation(tx, reservation); err != nil {
				continue
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
