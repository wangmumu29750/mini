package service

import (
	"math"
	"net/http"
	"strings"

	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"
)

const baseFareCentsPerKM = 13

type seatClassRule struct {
	Code        string
	Name        string
	Coefficient float64
}

var seatClassRulesByTrainType = map[string][]seatClassRule{
	"G": {
		{Code: "BUSINESS", Name: "商务座", Coefficient: 10.0},
		{Code: "FIRST", Name: "一等座", Coefficient: 5.8},
		{Code: "SECOND", Name: "二等座", Coefficient: 3.5},
	},
	"C": {
		{Code: "BUSINESS", Name: "商务座", Coefficient: 10.0},
		{Code: "FIRST", Name: "一等座", Coefficient: 5.8},
		{Code: "SECOND", Name: "二等座", Coefficient: 3.5},
	},
	"D": {
		{Code: "FIRST_SLEEPER", Name: "一等卧", Coefficient: 5.0},
		{Code: "SECOND_SLEEPER", Name: "二等卧", Coefficient: 3.8},
		{Code: "SECOND", Name: "二等座", Coefficient: 3.0},
	},
	"Z": conventionalSeatRules(),
	"T": conventionalSeatRules(),
	"K": conventionalSeatRules(),
}

func conventionalSeatRules() []seatClassRule {
	return []seatClassRule{
		{Code: "DELUXE_SOFT_SLEEPER", Name: "高级软卧", Coefficient: 4.0},
		{Code: "SOFT_SLEEPER", Name: "软卧", Coefficient: 3.0},
		{Code: "HARD_SLEEPER", Name: "硬卧", Coefficient: 2.0},
		{Code: "HARD_SEAT", Name: "硬座", Coefficient: 1.0},
		{Code: "NO_SEAT", Name: "无座", Coefficient: 1.0},
	}
}

func normalizeTrainType(value string) string {
	text := strings.ToUpper(strings.TrimSpace(value))
	if text == "" {
		return ""
	}
	return text[:1]
}

func supportedTrainType(value string) bool {
	_, ok := seatClassRulesByTrainType[normalizeTrainType(value)]
	return ok
}

func supportedSeatClasses(trainType string) []seatClassRule {
	return seatClassRulesByTrainType[normalizeTrainType(trainType)]
}

func trainTypeFromTrainNo(trainNo string) string {
	return normalizeTrainType(trainNo)
}

func ensureSeatClassAllowed(trainType, seatClassCode string) error {
	code := strings.ToUpper(strings.TrimSpace(seatClassCode))
	for _, rule := range supportedSeatClasses(trainType) {
		if rule.Code == code {
			return nil
		}
	}
	return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "席别不适用于该车次类型")
}

func calculateFareCents(mileage int, trainType, seatClassCode string) (int64, error) {
	if mileage <= 0 {
		return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "区间里程必须大于 0")
	}
	code := strings.ToUpper(strings.TrimSpace(seatClassCode))
	for _, rule := range supportedSeatClasses(trainType) {
		if rule.Code == code {
			price := float64(mileage) * baseFareCentsPerKM * rule.Coefficient
			return int64(math.Round(price/100) * 100), nil
		}
	}
	return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "席别不适用于该车次类型")
}

func seatClassName(code string) string {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	for _, rules := range seatClassRulesByTrainType {
		for _, rule := range rules {
			if rule.Code == normalized {
				return rule.Name
			}
		}
	}
	return code
}
