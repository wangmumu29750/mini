package response

import (
	"mini-12306/backend/pkg/requestctx"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK                    = "OK"
	CodeValidationError       = "VALIDATION_ERROR"
	CodeUnauthorized          = "UNAUTHORIZED"
	CodeForbidden             = "FORBIDDEN"
	CodeNotFound              = "NOT_FOUND"
	CodeConflict              = "CONFLICT"
	CodeInsufficientInventory = "INSUFFICIENT_INVENTORY"
	CodeInvalidOrderState     = "INVALID_ORDER_STATE"
	CodePaymentFailed         = "PAYMENT_FAILED"
	CodeTicketNotRefundable   = "TICKET_NOT_REFUNDABLE"
	CodeTicketNotChangeable   = "TICKET_NOT_CHANGEABLE"
	CodeInternalError         = "INTERNAL_ERROR"
)

type Body struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"traceId,omitempty"`
}

func OK(c *gin.Context, data any) {
	JSON(c, 200, CodeOK, "success", data)
}

func Error(c *gin.Context, httpStatus int, code string, message string) {
	JSON(c, httpStatus, code, message, nil)
}

func JSON(c *gin.Context, httpStatus int, code string, message string, data any) {
	traceID, _ := c.Get(requestctx.TraceIDKey)
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
		Data:    data,
		TraceID: stringValue(traceID),
	})
}

func stringValue(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}
