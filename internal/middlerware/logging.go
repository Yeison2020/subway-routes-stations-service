package middlerware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)


func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		start := time.Now()

		path := ctx.Request.URL.Path

		raw := ctx.Request.URL.RawQuery


		// Process request

		ctx.Next()

		// Log after processing

		latency := time.Since(start)

		clientIP := ctx.ClientIP()

		method := ctx.Request.Method

		statusCode := ctx.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency" : latency,
			"client_ip": clientIP,
			"method": method,
			"path": path,
			"user_agent": ctx.Request.UserAgent(),
		}).Info("HTTP Request")
	}

}