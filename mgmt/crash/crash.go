package crash

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler crashes the service
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "Crashing...")
		panic("Crashing...")
	}
}

// Plug plugs the API to /mgmt/crash
func Plug(r gin.IRoutes) {
	r.GET("/mgmt/crash", Handler())
}
