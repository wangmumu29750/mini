package service

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"mini-12306/backend/internal/model"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"
)

type TicketPriceContext struct {
	BasePriceCents int64
	SeatType       string
	TicketType     string
}

type TicketPriceRule interface {
	Apply(ctx TicketPriceContext) (int64, error)
}

type ticketPriceCalculator struct {
	rules map[string]TicketPriceRule
}

func newTicketPriceCalculator() *ticketPriceCalculator {
	return &ticketPriceCalculator{
		rules: map[string]TicketPriceRule{
			string(model.TicketTypeAdult):   adultTicketRule{},
			string(model.TicketTypeStudent): studentTicketRule{},
			string(model.TicketTypeChild):   childTicketRule{},
		},
	}
}

func CalculateTicketPrice(basePriceCents int64, seatType, ticketType string) (int64, error) {
	return newTicketPriceCalculator().Calculate(basePriceCents, seatType, ticketType)
}

func (c *ticketPriceCalculator) Calculate(basePriceCents int64, seatType, ticketType string) (int64, error) {
	seatType = strings.ToUpper(strings.TrimSpace(seatType))
	ticketType = strings.ToUpper(strings.TrimSpace(ticketType))
	if seatType == "" || ticketType == "" {
		return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "席别和票种不能为空")
	}
	rule, ok := c.rules[ticketType]
	if !ok {
		return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, fmt.Sprintf("不支持的票种: %s", ticketType))
	}
	return rule.Apply(TicketPriceContext{
		BasePriceCents: basePriceCents,
		SeatType:       seatType,
		TicketType:     ticketType,
	})
}

type adultTicketRule struct{}

func (adultTicketRule) Apply(ctx TicketPriceContext) (int64, error) {
	return ctx.BasePriceCents, nil
}

type studentTicketRule struct{}

func (studentTicketRule) Apply(ctx TicketPriceContext) (int64, error) {
	if ctx.SeatType != "SECOND" {
		return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "学生票仅支持二等座")
	}
	return roundedPrice(ctx.BasePriceCents, 0.75), nil
}

type childTicketRule struct{}

func (childTicketRule) Apply(ctx TicketPriceContext) (int64, error) {
	return roundedPrice(ctx.BasePriceCents, 0.5), nil
}

func roundedPrice(base int64, factor float64) int64 {
	if base <= 0 {
		return 0
	}
	return int64(math.Round(float64(base) * factor))
}
