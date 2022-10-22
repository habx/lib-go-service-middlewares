package querytoheader

import "github.com/gin-gonic/gin"

// Handler is a middleware that copies some query parameters to the headers
func Handler(conv map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range conv {
			if val := c.Query(k); val != "" {
				c.Request.Header.Set(v, val)
			}
		}

		c.Next()
	}
}
