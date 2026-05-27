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

func RoleRequired(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[role] = true
	}

	return func(c *gin.Context) {
		principal, ok := c.Get(requestctx.CurrentUserKey)
		if !ok {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "璇峰厛鐧诲綍")
			c.Abort()
			return
		}
		user, ok := principal.(auth.Principal)
		if !ok || !allowed[user.Role] {
			response.Error(c, http.StatusForbidden, response.CodeForbidden, "鏃犳潈闄?")
			c.Abort()
			return
		}
		c.Next()
	}
}
