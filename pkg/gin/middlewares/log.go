package middlewares

import (
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger is the logrus logger handler
func Logger(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()

		c.Next()

		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		entry := log.With(
			"host", c.Request.Host,
			"status", statusCode,
			"latency", latency, // time to process
			"ip", clientIP,
			"method", c.Request.Method,
			"length", c.Request.ContentLength,
			"ua", c.Request.UserAgent(),
			"uri", c.Request.RequestURI,
		)

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("\"%s %s\" %d (%dms)", c.Request.Method, path, statusCode, latency)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
