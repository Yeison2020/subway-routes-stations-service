package middlerware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)


const (
	 ErrCodeInternal = 500
	 
)

// PanicHandler catches panics
func PanicHandler() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        var details interface{}
        if err, ok := recovered.(string); ok {
            details = err
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    ErrCodeInternal,
                "message": "Internal server error",
                "details": details,
            },
        })
        c.Abort()
    })
}

// GlobalErrorHandler handles errors returned from handlers
func GlobalErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // process request

        errs := c.Errors
        if len(errs) > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": gin.H{
                    "code":    ErrCodeInternal,
                    "message": "Internal server error",
                    "details": errs[0].Error(),
                },
            })
        }
    }
}