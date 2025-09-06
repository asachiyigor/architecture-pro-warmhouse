package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin.HandlerFunc for logging requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("API Gateway - %v | %3d | %13v | %15s | %-7s %#v %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	})
}
