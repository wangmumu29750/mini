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
	TrainType      string
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

func CalculateTicketPrice(basePriceCents int64, trainType, seatType, ticketType string) (int64, error) {
	return newTicketPriceCalculator().Calculate(basePriceCents, trainType, seatType, ticketType)
}

func (c *ticketPriceCalculator) Calculate(basePriceCents int64, trainType, seatType, ticketType string) (int64, error) {
	trainType = strings.ToUpper(strings.TrimSpace(trainType))
	seatType = strings.ToUpper(strings.TrimSpace(seatType))
	ticketType = strings.ToUpper(strings.TrimSpace(ticketType))
	if trainType == "" || seatType == "" || ticketType == "" {
		return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "甯埆銆佽溅绫诲拰绁ㄧ涓嶈兘涓虹┖")
	}
	rule, ok := c.rules[ticketType]
	if !ok {
		return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, fmt.Sprintf("涓嶆敮鎸佺殑绁ㄧ: %s", ticketType))
	}
	return rule.Apply(TicketPriceContext{
		BasePriceCents: basePriceCents,
		TrainType:      trainType,
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
	switch ctx.TrainType {
	case "Z", "T", "K":
		return roundedPrice(ctx.BasePriceCents, 0.60), nil
	case "G", "C", "D":
		return roundedPrice(ctx.BasePriceCents, 0.75), nil
	default:
		return roundedPrice(ctx.BasePriceCents, 0.75), nil
	}
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
