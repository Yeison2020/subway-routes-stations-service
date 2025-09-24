package middlerware

import (
	"strconv"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"


	"github.com/yeison2020/subway-routing-service/internal/config"

	
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
	return func(ctx *gin.Context) {

		start := time.Now()
		ctx.Next() // process request
		duration := time.Since(start)

		cfg := config.LoadConfig()
		span, _ := tracer.SpanFromContext(ctx.Request.Context())
		if span == nil {
			log.Println("No active span found in the context; Logs might will have empty trace/spans ids")
		}


		// Prepare log fields
		fields := []zap.Field{
			zap.String("service", cfg.Service),
			zap.String("env", cfg.Env),
			zap.String("Content-Type", ctx.ContentType()),
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("user_agent", ctx.Request.UserAgent()),
			zap.String("client_ip", ctx.Request.Host),
			zap.Duration("latency", duration),
			zap.String("request_id", ctx.GetString("X-Request-ID")),
			zap.String("trace_id", span.Context().TraceID()),
			zap.String("span_id", strconv.FormatUint(span.Context().SpanID(), 10)),
		}

		// Choose level
		lvl := zap.InfoLevel
		if ctx.Writer.Status() >= 500 {
			lvl = zap.ErrorLevel
		} else if ctx.Writer.Status() >= 400 {
			lvl = zap.WarnLevel
		}

		// Log
		logger.Check(lvl, "HTTP request").Write(fields...)
	}
}

func RequestIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Try to get request ID from incoming header
		requestID := ctx.GetHeader("X-Request-ID")
		// If not present, generate a new one
		if requestID == "" {
			requestID = uuid.New().String()
		}
		// Set it back to headers so downstream systems see it
		ctx.Writer.Header().Set("X-Request-ID", requestID)

		// Store in Gin context for handlers/loggers
		ctx.Set("X-Request-ID", requestID)

		// Continue
		ctx.Next()

	}
}
