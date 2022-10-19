package mgmt

import (
	"github.com/gin-gonic/gin"

	buildt "github.com/habx/lib-go-types/build"
	"github.com/habx/lib-go-utils/health"
)

// Plug the mgmt API to the given router
func Plug(eng *gin.Engine, options ...PlugOption) {
	conf := &plugConfig{}

	for _, opt := range options {
		opt(conf)
	}

	r := eng.Group("/mgmt")
	r.GET("/version", VersionHandler(conf.buildInfo))
	r.GET("/health", HealthHandler(conf.healthManager))
	r.GET("/crash", CrashHandler())
}

type plugConfig struct {
	healthManager *health.Manager
	buildInfo     *buildt.Info
}

// PlugOption is a function that configures the plug
type PlugOption func(*plugConfig)

// OptHealthManager sets the health manager to use
func OptHealthManager(m *health.Manager) PlugOption {
	return func(c *plugConfig) {
		c.healthManager = m
	}
}

// OptBuildInfo sets the build info to use
func OptBuildInfo(bi *buildt.Info) PlugOption {
	return func(c *plugConfig) {
		c.buildInfo = bi
	}
}
