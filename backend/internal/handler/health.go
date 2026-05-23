package handler

import (
	"net/http"
	"time"

	"mini-12306/backend/internal/config"
	"mini-12306/backend/internal/database"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	cfg config.Config
	db  *gorm.DB
}

type HealthResponse struct {
	Status   string `json:"status"`
	Env      string `json:"env"`
	Database string `json:"database"`
	Time     string `json:"time"`
}

func NewHealthHandler(cfg config.Config, db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		cfg: cfg,
		db:  db,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	statusCode := http.StatusOK
	code := response.CodeOK
	message := "success"
	status := "ok"
	databaseStatus := "ok"

	if err := database.Ping(h.db); err != nil {
		statusCode = http.StatusServiceUnavailable
		code = response.CodeInternalError
		message = "数据库暂时不可用"
		status = "degraded"
		databaseStatus = "unavailable"
	}

	response.JSON(c, statusCode, code, message, HealthResponse{
		Status:   status,
		Env:      h.cfg.AppEnv,
		Database: databaseStatus,
		Time:     time.Now().Format(time.RFC3339),
	})
}
