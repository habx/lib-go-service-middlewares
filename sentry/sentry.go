package sentry

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func Handler() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{Repanic: true})
}
