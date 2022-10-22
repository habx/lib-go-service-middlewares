package version

import (
	"net/http"

	"github.com/gin-gonic/gin"

	buildt "github.com/habx/lib-go-types/build"
	buildu "github.com/habx/lib-go-utils/build"
)

// Handler returns the version of the service
func Handler(bi *buildt.Info) gin.HandlerFunc {
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

// Plug plugs the API to /mgmt/version
func Plug(r gin.IRoutes, bi *buildt.Info) {
	r.GET("/mgmt/version", Handler(bi))
}
