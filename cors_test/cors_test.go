package cors_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	tgin "github.com/habx/lib-go-tests/http/gin"

	"github.com/habx/lib-go-service-middlewares/cors"
)

func TestCors(t *testing.T) {
	srv, eng := tgin.GetServerWithGin(t)

	eng.POST("/test", cors.Handler(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	createReq := func(t *testing.T) *http.Request {
		req, err := http.NewRequest(http.MethodPost, srv.URL("/test?name"+t.Name()), nil)
		assert.NoError(t, err)

		return req
	}

	check := func(t *testing.T, req *http.Request, good bool) {
		a := assert.New(t)

		resp, err := http.DefaultClient.Do(req)
		a.NoError(err)

		defer resp.Body.Close()

		if good {
			a.Equal(http.StatusOK, resp.StatusCode)
		} else {
			a.Equal(http.StatusForbidden, resp.StatusCode)
		}
	}

	t.Run("good cors", func(t *testing.T) {
		req := createReq(t)
		req.Header["Origin"] = []string{"https://www.habx.com"}
		check(t, req, true)
	})

	t.Run("bad cors origin", func(t *testing.T) {
		req := createReq(t)
		req.Header["Origin"] = []string{"https://www.google.com"}
		check(t, req, false)
	})

	t.Run("bad cors header", func(t *testing.T) {
		t.Skip()
		req := createReq(t)
		req.Header["Origin"] = []string{"https://www.habx.com"}
		req.Header["Something"] = []string{"bad ?"}
		check(t, req, false)
	})
}
