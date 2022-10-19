package mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/habx/lib-go-utils/health"
)

// HealthHandler returns the version of the service
func HealthHandler(m *health.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := m.Check(); err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.String(http.StatusOK, "OK")
	}
}

// HealthPlug plugs the API to /mgmt/version
func HealthPlug(eng *gin.Engine, m *health.Manager) {
	eng.GET("/mgmt/health", HealthHandler(m))
}
