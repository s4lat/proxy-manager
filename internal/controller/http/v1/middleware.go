package v1

import (
	"encoding/json"
	"proxy_manager/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// JSONLogMiddleware logs a gin HTTP request in JSON format, with some additional custom key/values.
func JSONLogMiddleware(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := time.Since(start).Milliseconds()

		logMsg, err := json.Marshal(map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.RequestURI,
			"duration": duration,
			"status":   c.Writer.Status(),
		})
		if err != nil {
			l.Error("JSONLogMiddleware: can't marshal log - %s", err.Error())
		} else {
			l.Info(string(logMsg))
		}
	}
}
