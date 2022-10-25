package version

import (
	"log"
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
			log.Println("WARNING: lib-go-service-middlewares: Couldn't fetch build info: ", err)
		}
	}

	return func(c *gin.Context) {
		if bi != nil {
			c.JSON(http.StatusOK, bi)
		} else {
			c.String(http.StatusInternalServerError, "No build info")
		}
	}
}

// Plug plugs the API to /mgmt/version
func Plug(r gin.IRoutes, bi *buildt.Info) {
	r.GET("/mgmt/version", Handler(bi))
}
