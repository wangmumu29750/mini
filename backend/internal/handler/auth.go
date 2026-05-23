package handler

import (
	"net/http"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/service"
	pkgauth "mini-12306/backend/pkg/auth"
	"mini-12306/backend/pkg/requestctx"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "注册信息不完整或格式不正确")
		return
	}

	result, err := h.auth.Register(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "用户名和密码不能为空")
		return
	}

	result, err := h.auth.Login(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	principal, ok := c.Get(requestctx.CurrentUserKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "请先登录")
		return
	}

	user, ok := principal.(pkgauth.Principal)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "登录状态已失效")
		return
	}

	result, err := h.auth.CurrentUser(user.UserID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.OK(c, nil)
}
