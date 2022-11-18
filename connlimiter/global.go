package connlimiter

import "github.com/gin-gonic/gin"

// From: https://github.com/aviddiviner/gin-limit/blob/master/limit.go
func Global(n int) gin.HandlerFunc {
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
