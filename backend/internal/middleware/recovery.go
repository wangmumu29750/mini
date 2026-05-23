package middleware

import (
	"log"
	"net/http"

	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Printf("panic recovered: %v", recovered)
		response.JSON(c, http.StatusInternalServerError, response.CodeInternalError, "服务暂时不可用", nil)
	})
}
