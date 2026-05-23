package handler

import (
	"errors"
	"log"
	"net/http"

	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.Status, appErr.Code, appErr.Message)
		return
	}

	log.Printf("request failed: %v", err)
	response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "服务暂时不可用")
}
