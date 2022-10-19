package mgmt

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	thttp "github.com/habx/lib-go-tests/http"
	buildt "github.com/habx/lib-go-types/build"

	"github.com/habx/lib-go-service-middlewares/mgmt"
)

func TestPlug(t *testing.T) {
	a := assert.New(t)

	eng := gin.New()
	mgmt.Plug(eng)

	srv := thttp.GetServer(t, thttp.OptHandler(eng))
	c := srv.GetClient()

	t.Run("version", func(t *testing.T) {
		a.Contains(c.GETAsString("/mgmt/version"), "\"version\":\"")
	})

	t.Run("crash", func(t *testing.T) {
		req, err := http.NewRequest("GET", srv.URL("/mgmt/crash"), nil)
		a.NoError(err)

		resp, err := c.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		a.Error(err)
		a.Nil(resp)
	})
}

func TestVersion(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	mgmt.VersionPlug(r, &buildt.Info{Version: "1.2.3"})

	s := thttp.GetServer(t, thttp.OptHandler(r))
	c := s.GetClient()
	a.Contains(c.GETAsString("/mgmt/version"), "1.2.3")
}
