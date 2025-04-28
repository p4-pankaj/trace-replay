package tracer

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceID := uuid.New().String()
		recorder := NewTraceRecorder(traceID)

		ctx := ToContext(c.Request.Context(), recorder)
		c.Request = c.Request.WithContext(ctx)

		recorder.Record(fmt.Sprintf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path))

		c.Next()

		recorder.Record(fmt.Sprintf("Response status: %d", c.Writer.Status()))
	}
}
