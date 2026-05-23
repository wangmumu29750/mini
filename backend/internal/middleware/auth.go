package middleware

import (
	"net/http"
	"strings"

	"mini-12306/backend/pkg/auth"
	"mini-12306/backend/pkg/requestctx"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		principal, err := auth.ParseToken(jwtSecret, strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "登录状态已失效")
			c.Abort()
			return
		}

		c.Set(requestctx.CurrentUserKey, principal)
		c.Next()
	}
}
