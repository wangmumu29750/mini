package handler

import (
	"net/http"
	"strconv"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/service"
	pkgauth "mini-12306/backend/pkg/auth"
	"mini-12306/backend/pkg/requestctx"
	"mini-12306/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orders *service.OrderService
}

func NewOrderHandler(orders *service.OrderService) *OrderHandler {
	return &OrderHandler{orders: orders}
}

func (h *OrderHandler) Create(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "下单信息不完整")
		return
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = c.GetHeader("Idempotency-Key")
	}

	result, err := h.orders.Create(principal.UserID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *OrderHandler) ClerkCreate(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	var req dto.ClerkCreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "售票下单信息不完整")
		return
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = c.GetHeader("Idempotency-Key")
	}

	result, err := h.orders.ClerkCreate(principal.UserID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *OrderHandler) List(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	result, err := h.orders.List(principal.UserID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *OrderHandler) Detail(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "订单编号不正确")
		return
	}

	result, err := h.orders.Detail(principal.UserID, orderID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *OrderHandler) Pay(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "订单编号不正确")
		return
	}

	var req dto.PayOrderRequest
	_ = c.ShouldBindJSON(&req)

	result, err := h.orders.Pay(principal.UserID, orderID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	principal, ok := currentPrincipal(c)
	if !ok {
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "订单编号不正确")
		return
	}

	result, err := h.orders.Cancel(principal.UserID, orderID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, result)
}

func currentPrincipal(c *gin.Context) (pkgauth.Principal, bool) {
	principal, ok := c.Get(requestctx.CurrentUserKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "请先登录")
		return pkgauth.Principal{}, false
	}
	user, ok := principal.(pkgauth.Principal)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "登录状态已失效")
		return pkgauth.Principal{}, false
	}
	return user, true
}
