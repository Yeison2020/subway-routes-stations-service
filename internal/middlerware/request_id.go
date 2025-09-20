package middlerware


import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func RequestIDMiddlerware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		
		requestID := ctx.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Header("X-Request-ID", requestID)
		ctx.Set("request_id", requestID)
		ctx.Next()

	}
}