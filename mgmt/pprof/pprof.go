package pprof

import (
	"errors"

	gpprof "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// ErrNoManager is returned when no ACL manager is provided
var ErrNoManager = errors.New("no manager provided")

// PlugOnRoute plugs the API to /mgmt/pprof
// Note: We're using a RouterGroup instead of IRoutes because gin's pprof module
// requires it. Not for very real reasons.
func PlugOnRoute(grp *gin.RouterGroup, route string) {
	gpprof.RouteRegister(grp, route)
}

// Plug plugs the API to /mgmt/pprof
func Plug(r gin.IRoutes) {
	PlugOnRoute(r.(*gin.RouterGroup), "/mgmt/pprof")
}
