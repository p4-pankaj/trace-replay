package tracer

import (
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/p4-pankaj/trace-replay/config"
	"github.com/p4-pankaj/trace-replay/db"
)

func TraceMiddleware(db db.TraceStorage, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		recorder := NewTraceRecorder(traceID, db)

		ctx := ToContext(c.Request.Context(), recorder)
		c.Request = c.Request.WithContext(ctx)
		recorder.Lg.Info().Msg(fmt.Sprintf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path))

		c.Next()
		fmt.Println("Calling CloseUpdater")
		go func() {
			<-c.Request.Context().Done()
			// after request is completed , give trace 29
			time.Sleep(2 * time.Second)
			recorder.CloseUpdater <- struct{}{}

		}()

	}
}
