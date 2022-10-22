package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/habx/lib-go-utils/health"
)

// Handler returns the version of the service
func Handler(m *health.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := m.Check(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// Plug plugs the API to /mgmt/version
func Plug(r gin.IRoutes, m *health.Manager) {
	r.GET("/mgmt/health", Handler(m))
}
