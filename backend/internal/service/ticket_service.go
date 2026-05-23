package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
)

type TicketService struct {
	tickets *repository.TicketRepository
}

func NewTicketService(tickets *repository.TicketRepository) *TicketService {
	return &TicketService{tickets: tickets}
}

func (s *TicketService) List(userID uint64) ([]dto.TicketResponse, error) {
	tickets, err := s.tickets.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TicketResponse, 0, len(tickets))
	for _, ticket := range tickets {
		result = append(result, ticketResponse(ticket))
	}
	return result, nil
}

func (s *TicketService) Detail(userID, ticketID uint64) (dto.TicketResponse, error) {
	ticket, err := s.tickets.FindByUserAndID(userID, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TicketResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "车票不存在")
		}
		return dto.TicketResponse{}, err
	}
	return ticketResponse(ticket), nil
}

func ticketResponse(ticket model.Ticket) dto.TicketResponse {
	return dto.TicketResponse{
		ID:             ticket.ID,
		TicketNo:       ticket.TicketNo,
		OrderID:        ticket.OrderID,
		TrainID:        ticket.TrainID,
		TrainNo:        ticket.TrainNo,
		TravelDate:     ticket.TravelDate.Format("2006-01-02"),
		FromStation:    dto.StationResponse{ID: ticket.FromStationID, Name: ticket.FromStationName},
		ToStation:      dto.StationResponse{ID: ticket.ToStationID, Name: ticket.ToStationName},
		SeatClassCode:  ticket.SeatClassCode,
		SeatClassName:  ticket.SeatClassName,
		PassengerName:  ticket.PassengerName,
		IDCardNoMasked: maskIDCardNo(ticket.IDCardNo),
		Status:         string(ticket.Status),
		IssuedAt:       ticket.IssuedAt.Format(time.RFC3339),
	}
}

func maskIDCardNo(value string) string {
	text := strings.TrimSpace(value)
	if len(text) <= 8 {
		return text
	}
	return text[:4] + strings.Repeat("*", len(text)-8) + text[len(text)-4:]
}
