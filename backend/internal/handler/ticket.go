package handler

import (
	"net/http"
	"strconv"

	"mini-12306/backend/internal/dto"
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

func (h *TicketHandler) Refund(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || ticketID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "车票编号不正确")
		return
	}

	var req dto.RefundTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "退票参数不正确")
		return
	}

	result, err := h.tickets.Refund(principal.UserID, ticketID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *TicketHandler) ChangeOptions(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || ticketID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "车票编号不正确")
		return
	}

	var query dto.ChangeOptionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "改签查询参数不正确")
		return
	}

	result, err := h.tickets.ChangeOptions(principal.UserID, ticketID, query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *TicketHandler) Change(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || ticketID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "车票编号不正确")
		return
	}

	var req dto.ChangeTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "改签参数不正确")
		return
	}

	result, err := h.tickets.Change(principal.UserID, ticketID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}
