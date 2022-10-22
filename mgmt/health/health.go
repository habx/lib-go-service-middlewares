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
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.String(http.StatusOK, "OK")
	}
}

// Plug plugs the API to /mgmt/version
func Plug(r gin.IRoutes, m *health.Manager) {
	r.GET("/mgmt/health", Handler(m))
}
