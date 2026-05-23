package repository

import (
	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *TicketRepository) ListByUser(userID uint64) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Where("user_id = ?", userID).
		Order("travel_date DESC, issued_at DESC").
		Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	if err := attachTicketTimes(r.db, tickets); err != nil {
		return nil, err
	}
	return tickets, err
}

func (r *TicketRepository) FindByUserAndID(userID, ticketID uint64) (model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket).Error
	if err != nil {
		return ticket, err
	}
	tickets := []model.Ticket{ticket}
	if err := attachTicketTimes(r.db, tickets); err != nil {
		return ticket, err
	}
	ticket = tickets[0]
	return ticket, err
}

func attachTicketTimes(db *gorm.DB, tickets []model.Ticket) error {
	for i := range tickets {
		depart, arrive, err := stationTimes(db, tickets[i].TravelDate, tickets[i].TrainID, tickets[i].FromStationID, tickets[i].ToStationID)
		if err != nil {
			return err
		}
		tickets[i].DepartTime = &depart
		tickets[i].ArriveTime = &arrive
	}
	return nil
}
