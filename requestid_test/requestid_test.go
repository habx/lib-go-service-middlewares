package requestidtest_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	thttp "github.com/habx/lib-go-tests/http"
	tgin "github.com/habx/lib-go-tests/http/gin"

	"github.com/habx/lib-go-service-middlewares/requestid"
)

func TestRequestId(t *testing.T) {
	srv, g := tgin.GetServerWithGin(t)

	r := g.Group("/", requestid.Handler())

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, requestid.Get(c))
	})

	t.Run("without header", func(t *testing.T) {
		a := assert.New(t)
		c := srv.GetClient()

		requestID := c.GetString("/test")
		a.NotEmpty(requestID)
		_, err := uuid.Parse(requestID)
		a.NoError(err)
	})

	t.Run("with header", func(t *testing.T) {
		a := assert.New(t)
		c := srv.GetClient(thttp.OptCltHeader("X-Request-ID", "something"))
		a.Equal("something", c.GetString("/test"))
	})
}
