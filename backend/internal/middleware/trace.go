package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"mini-12306/backend/pkg/requestctx"

	"github.com/gin-gonic/gin"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = newTraceID()
		}

		c.Set(requestctx.TraceIDKey, traceID)
		c.Header("X-Trace-Id", traceID)
		c.Next()
	}
}

func newTraceID() string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return "req_" + hex.EncodeToString(bytes[:])
	}
	return "req_" + hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
}
