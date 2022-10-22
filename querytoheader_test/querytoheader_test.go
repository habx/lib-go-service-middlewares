package mgmt_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	thttp "github.com/habx/lib-go-tests/http"

	"github.com/habx/lib-go-service-middlewares/querytoheader"
)

func TestQueryToHeader(t *testing.T) {
	a := assert.New(t)

	g := gin.Default()

	r := g.Group("/", querytoheader.Handler(map[string]string{"param": "header"}))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, c.Request.Header.Get("header"))
	})

	srv := thttp.GetServer(t, thttp.OptHandler(g))
	c := srv.GetClient()
	a.Equal("foo", c.GetString("/test?param=foo"))
}
