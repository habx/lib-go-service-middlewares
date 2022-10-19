package mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"

	buildt "github.com/habx/lib-go-types/build"
	buildu "github.com/habx/lib-go-utils/build"
)

// VersionHandler returns the version of the service
func VersionHandler(bi *buildt.Info) gin.HandlerFunc {
	if bi == nil {
		var err error
		bi, err = buildu.GetInfo()

		if err != nil {
			panic(err)
		}
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, bi)
	}
}

// VersionPlug plugs the API to /mgmt/version
func VersionPlug(eng *gin.Engine, bi *buildt.Info) {
	eng.GET("/mgmt/version", VersionHandler(bi))
}
