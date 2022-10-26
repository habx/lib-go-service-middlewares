package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Header is the header name used to store the request id
const header = "X-Request-Id"
const contextKey = "requestId"

// Handler returns a gin middleware that adds a request id to the context
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(header)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header(header, requestID)
		}

		c.Set(contextKey, requestID)
		c.Next()
	}
}

// Get returns the request id from the context
func Get(c *gin.Context) string {
	return c.GetString(contextKey)
}
