package memstats

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/habx/lib-go-utils/cgstats"
)

type memStats struct {
	Runtime runtime.MemStats  `json:"runtime"`
	CGroup  map[string]uint64 `json:"cgroup"`
	Status  string            `json:"status,omitempty"`
}

// Handler returns the memory stats of the service
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m memStats

		runtime.ReadMemStats(&m.Runtime)

		var err error
		m.CGroup, err = cgstats.GetAllMemoryStats()

		if err != nil {
			m.Status += "cgstats:" + err.Error() + "\n"
		}

		c.JSON(http.StatusOK, m)
	}
}

// Plug plugs the API to /mgmt/memstats
func Plug(r gin.IRoutes) {
	r.GET("/mgmt/memstats", Handler())
}
