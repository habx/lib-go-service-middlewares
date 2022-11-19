package connlimiter

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// QueueGlobal is a global connection limiter
// From: https://github.com/aviddiviner/gin-limit/blob/master/limit.go
func QueueGlobal(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(c *gin.Context) {
		acquire()
		// before request
		defer release() // after request
		c.Next()
	}
}

// DropGlobal limits the number of concurrent connections by dropping any new connection after that limit
func DropGlobal(n int) gin.HandlerFunc {
	nbConnections := int32(0)
	maxNbConnections := int32(n)

	return func(c *gin.Context) {
		atomic.AddInt32(&nbConnections, 1)
		defer atomic.AddInt32(&nbConnections, -1)

		if nbConnections > maxNbConnections {
			c.AbortWithStatus(http.StatusTooManyRequests)

			return
		}

		c.Next()
	}
}
