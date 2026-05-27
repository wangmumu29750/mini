package handler

import (
	"net/http"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/service"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemSettingHandler struct {
	settings *service.SystemSettingService
}

func NewSystemSettingHandler(settings *service.SystemSettingService) *SystemSettingHandler {
	return &SystemSettingHandler{settings: settings}
}

func (h *SystemSettingHandler) List(c *gin.Context) {
	result, err := h.settings.List()
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *SystemSettingHandler) Update(c *gin.Context) {
	var req dto.UpdateSystemSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "system setting payload is invalid")
		return
	}

	result, err := h.settings.Update(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}
