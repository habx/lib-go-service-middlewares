package cors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	tgin "github.com/habx/lib-go-tests/http/gin"

	"github.com/habx/lib-go-service-middlewares/cors"
)

func checkResponse(t *testing.T, req *http.Request, good bool) {
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

func TestCors(t *testing.T) {
	srv, eng := tgin.GetServerWithGin(t)

	eng.POST("/test", cors.Handler(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	eng.POST("/test2", cors.Handler(cors.OptForceTransmission()), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	createReq := func(t *testing.T, forceTransmission bool) *http.Request {
		var path string
		if forceTransmission {
			path = "test2"
		} else {
			path = "test"
		}

		req, err := http.NewRequest(http.MethodPost, srv.URL(fmt.Sprintf("/%s?name%s", path, t.Name())), nil)
		assert.NoError(t, err)

		return req
	}

	t.Run("good cors", func(t *testing.T) {
		req := createReq(t, false)
		req.Header["Origin"] = []string{"https://www.habx.com"}
		checkResponse(t, req, true)
	})

	t.Run("good cors (force)", func(t *testing.T) {
		req := createReq(t, true)
		req.Header["Origin"] = []string{"https://www.habx.com"}
		checkResponse(t, req, true)
	})

	t.Run("bad cors origin", func(t *testing.T) {
		req := createReq(t, false)
		req.Header["Origin"] = []string{"https://www.google.com"}
		checkResponse(t, req, false)
	})

	t.Run("bad cors header", func(t *testing.T) {
		t.Skip()
		req := createReq(t, false)
		req.Header["Origin"] = []string{"https://www.habx.com"}
		req.Header["Something"] = []string{"bad ?"}
		checkResponse(t, req, false)
	})
}
