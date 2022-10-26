package sentry_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	tgin "github.com/habx/lib-go-tests/http/gin"

	"github.com/habx/lib-go-service-middlewares/sentry"
)

func TestSentry(t *testing.T) {
	t.Run("without recovery", func(t *testing.T) {
		a := assert.New(t)
		a.Panics(func() {
			srv, g := tgin.GetServerWithGin(t)

			r := g.Group("/", sentry.Handler())

			r.GET("/crash", func(c *gin.Context) {
				panic("crashing")
			})

			c := srv.GetClient()
			a.Equal(http.StatusInternalServerError, c.GetStatusCode("/crash"))
		})
	})

	t.Run("with recovery", func(t *testing.T) {
		a := assert.New(t)

		srv, g := tgin.GetServerWithGin(t)

		r := g.Group("/", gin.Recovery(), sentry.Handler())

		r.GET("/crash", func(c *gin.Context) {
			panic("crashing")
		})

		c := srv.GetClient()
		a.Equal(http.StatusInternalServerError, c.GetStatusCode("/crash"))
	})
}
