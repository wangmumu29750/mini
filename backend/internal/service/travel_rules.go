package service

import (
	"net/http"
	"time"

	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"
)

func validateTicketTravelDate(date time.Time) error {
	if startOfLocalDay(date).Before(startOfLocalDay(time.Now())) {
		return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "过期日期不能购买车票")
	}
	return nil
}

func validateFutureDeparture(travelDate time.Time, departClock string, dayOffset int) error {
	if !combineDateClock(travelDate, departClock, dayOffset).After(time.Now()) {
		return apperrors.New(http.StatusConflict, response.CodeValidationError, "当前车次发车时间已过，不能购买或改签到该车次")
	}
	return nil
}

func percentageFee(amountCents int64, percent int) int64 {
	if amountCents <= 0 || percent <= 0 {
		return 0
	}
	return amountCents * int64(percent) / 100
}
