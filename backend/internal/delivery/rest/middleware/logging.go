package middleware

import (
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// sensitiveQueryKeys lists query parameter names whose values should be redacted in logs.
var sensitiveQueryKeys = map[string]bool{
	"api_key":  true,
	"token":    true,
	"secret":   true,
	"password": true,
}

// sanitizeQuery parses a raw query string and replaces the values of
// sensitive parameters with "[REDACTED]".
func sanitizeQuery(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}

	params, err := url.ParseQuery(rawQuery)
	if err != nil {
		// If we can't parse it, redact entirely to be safe.
		return "[REDACTED]"
	}

	for key := range params {
		if sensitiveQueryKeys[strings.ToLower(key)] {
			params.Set(key, "[REDACTED]")
		}
	}

	return params.Encode()
}

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := sanitizeQuery(c.Request.URL.RawQuery)

		c.Next()

		duration := time.Since(start)

		logger.Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
