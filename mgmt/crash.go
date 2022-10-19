package mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CrashHandler crashes the service
func CrashHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "Crashing...")
		panic("Crashing...")
	}
}

// CrashPlug plugs the API to /mgmt/crash
func CrashPlug(eng *gin.Engine) {
	eng.GET("/mgmt/crash", CrashHandler())
}
