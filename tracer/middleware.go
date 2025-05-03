package tracer

import (
	"fmt"

	"github.com/p4-pankaj/trace-replay/db"
	"github.com/p4-pankaj/trace-replay/traceConfig"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware(db db.TraceStorage, config *traceConfig.TraceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var recorder *TraceRecorder
		if config.Env == "DEBUG" {
			recorder = NewTraceRecorder(c, config.
				DebugConfig.TraceId, db, config.Env)
		} else {
			recorder = NewTraceRecorder(c, uuid.New().String(),
				db, config.Env)
		}
		ctx := ToContext(c.Request.Context(), recorder)
		c.Request = c.Request.WithContext(ctx)
		recorder.Lg.Info().Msg(fmt.Sprintf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path))
		c.Set("tracer", recorder)
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
