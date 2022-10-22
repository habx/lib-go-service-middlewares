package memstats

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Handler returns the memory stats of the service
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m runtime.MemStats

		runtime.ReadMemStats(&m)

		c.JSON(http.StatusOK, m)
	}
}

// Plug plugs the API to /mgmt/memstats
func Plug(r gin.IRoutes) {
	r.GET("/mgmt/memstats", Handler())
}
