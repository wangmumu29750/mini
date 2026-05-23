package handler

import (
	"net/http"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/service"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type TrainHandler struct {
	trains *service.TrainService
}

func NewTrainHandler(trains *service.TrainService) *TrainHandler {
	return &TrainHandler{trains: trains}
}

func (h *TrainHandler) Stations(c *gin.Context) {
	result, err := h.trains.ListStations()
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *TrainHandler) Search(c *gin.Context) {
	var query dto.TrainSearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "请填写乘车日期、出发站和到达站")
		return
	}

	result, err := h.trains.Search(query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}
