package handler

import (
	"net/http"
	"strconv"

	"mini-12306/backend/internal/service"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	tickets *service.TicketService
}

func NewTicketHandler(tickets *service.TicketService) *TicketHandler {
	return &TicketHandler{tickets: tickets}
}

func (h *TicketHandler) List(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	result, err := h.tickets.List(principal.UserID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *TicketHandler) Detail(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || ticketID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "车票编号不正确")
		return
	}

	result, err := h.tickets.Detail(principal.UserID, ticketID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}
