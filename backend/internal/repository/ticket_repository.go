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

func (r *TicketRepository) ListByUser(userID uint64) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Where("user_id = ?", userID).
		Order("travel_date DESC, issued_at DESC").
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) FindByUserAndID(userID, ticketID uint64) (model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket).Error
	return ticket, err
}
