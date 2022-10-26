package sentry

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// Handler returns a gin middleware that sends errors to sentry
func Handler() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{Repanic: true})
}
