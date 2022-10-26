package mgmt_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	tgin "github.com/habx/lib-go-tests/http/gin"

	"github.com/habx/lib-go-service-middlewares/querytoheader"
)

func TestQueryToHeader(t *testing.T) {
	a := assert.New(t)

	srv, g := tgin.GetServerWithGin(t)

	r := g.Group("/", querytoheader.Handler(map[string]string{"param": "header"}))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, c.Request.Header.Get("header"))
	})

	c := srv.GetClient()
	a.Equal("foo", c.GetString("/test?param=foo"))
}
