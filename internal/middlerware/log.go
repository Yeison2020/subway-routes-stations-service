package middlerware

import (

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var Logger *zap.Logger

func InitLogger() {
	// Create a development config
	cfg := zap.NewDevelopmentConfig()

	// Disable stack traces and caller info
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}

}


// LoggerMiddleware logs requests with latency, status, method, and path
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // process request
		duration := time.Since(start)


		// Prepare log fields
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("client_ip", c.Request.Host),
			zap.Duration("latency", duration),
		}

		// Choose level
		lvl := zap.InfoLevel
		if c.Writer.Status() >= 500 {
			lvl = zap.ErrorLevel
		} else if c.Writer.Status() >= 400 {
			lvl = zap.WarnLevel
		}

		// Log
		logger.Check(lvl, "HTTP request").Write(fields...)
	}
}
