package middlerware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ErrCodeInternal = 500
)

// PanicHandler catches panics
func PanicHandler(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
	logger.Error("Panic recovered",
            zap.Any("panic", recovered),
            zap.ByteString("stack", debug.Stack()),
        )

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    ErrCodeInternal,
				"message": "Internal server error",
			},
		})
		ctx.Abort()
	})
}
