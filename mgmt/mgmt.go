package mgmt

import (
	"github.com/gin-gonic/gin"
)

// Plug the mgmt API to the given router
func Plug(eng *gin.Engine) {
	r := eng.Group("/mgmt")
	r.GET("/version", VersionHandler(nil))
	r.GET("/crash", CrashHandler())
}
