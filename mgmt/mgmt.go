// Package mgmt provides a set of management endpoints. If you don't want to
// import a bunch of dependencies, please use the subpackages directly.
package mgmt

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	acl "github.com/habx/lib-go-acl/v2"
	aclmid "github.com/habx/lib-go-acl/v2/middleware"
	tauth "github.com/habx/lib-go-types/auth"
	tbuild "github.com/habx/lib-go-types/build"
	uhealth "github.com/habx/lib-go-utils/health"

	"github.com/habx/lib-go-service-middlewares/mgmt/crash"
	"github.com/habx/lib-go-service-middlewares/mgmt/health"
	"github.com/habx/lib-go-service-middlewares/mgmt/memstats"
	"github.com/habx/lib-go-service-middlewares/mgmt/pprof"
	"github.com/habx/lib-go-service-middlewares/mgmt/version"
	"github.com/habx/lib-go-service-middlewares/querytoheader"
)

const aclSlug = "tech.diagnostics.access"

// Plug the mgmt API to the given router
func Plug(eng *gin.Engine, options ...PlugOption) error {
	conf, err := initConfig(options...)

	if err != nil {
		return err
	}

	mid, err := aclmid.NewAuthMidleware(
		aclmid.OptManager(conf.aclManager),
		aclmid.OptAccessACL(aclSlug),
	)

	if err != nil {
		return fmt.Errorf("unable to create ACL middleware: %w", err)
	}

	r := eng.Group("/mgmt")

	// These calls are without auth:
	r.HEAD("/version")
	r.GET("/version", version.Handler(conf.buildInfo))
	r.HEAD("/health")
	r.GET("/health", health.Handler(conf.healthManager))
	r.GET("/crash", crash.Handler())

	// These calls are with auth:
	r.GET("/memstats", mid, memstats.Handler())
	pprof.PlugOnRoute(r.Group("/pprof", querytoheader.Handler(map[string]string{"token": "auth-token"}), mid), "/")

	if conf.token.IsSet() {
		_, err = conf.aclManager.DeclarePermissions(
			[]acl.PermissionInput{{
				Slug:         aclSlug,
				Name:         "Access to the diagnostics endpoints",
				InitialRoles: []string{"admin"},
			}},
			conf.token,
		)

		// We can't just fail because of an expired token
		if err != nil {
			conf.logger.Warnw("unable to declare permissions", "err", err)
		}
	}

	return nil
}

type plugConfig struct {
	healthManager *uhealth.Manager
	aclManager    *acl.Manager
	buildInfo     *tbuild.Info
	logger        *zap.SugaredLogger
	token         tauth.Token
	habxEnv       string
}

// ErrMisingHabxEnv is returned when the ACL manager is not set and the habx env is not set
var ErrMisingHabxEnv = fmt.Errorf("missing habxEnv")

func initConfig(options ...PlugOption) (*plugConfig, error) {
	conf := &plugConfig{}

	for _, opt := range options {
		opt(conf)
	}

	// It's best to pass the habxEnv through the OptHabxEnv option, but to make
	// this API convenient we can skip it.
	if conf.habxEnv == "" {
		conf.habxEnv = os.Getenv("HABX_ENV")
	}

	if conf.habxEnv == "" && conf.aclManager == nil {
		return nil, ErrMisingHabxEnv
	}

	if conf.logger == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("unable to create logger: %w", err)
		}

		conf.logger = logger.Sugar()
	}

	if conf.aclManager == nil {
		var err error
		conf.aclManager, err = acl.NewACLManager(acl.OptHabxEnv(conf.habxEnv), acl.OptLogger(conf.logger))

		if err != nil {
			return nil, fmt.Errorf("unable to create ACL manager: %w", err)
		}
	}

	if conf.healthManager == nil {
		conf.healthManager = uhealth.NewManager()
	}

	return conf, nil
}

// PlugOption is a function that configures the plug
type PlugOption func(*plugConfig)

// OptHealthManager sets the health manager to use
func OptHealthManager(m *uhealth.Manager) PlugOption {
	return func(c *plugConfig) {
		c.healthManager = m
	}
}

// OptACLManager sets the ACL manager to use
func OptACLManager(m *acl.Manager) PlugOption {
	return func(c *plugConfig) {
		c.aclManager = m
	}
}

// OptBuildInfo sets the build info to use
func OptBuildInfo(bi *tbuild.Info) PlugOption {
	return func(c *plugConfig) {
		c.buildInfo = bi
	}
}

// OptToken sets the token to use
func OptToken(t tauth.Token) PlugOption {
	return func(c *plugConfig) {
		c.token = t
	}
}

// OptLogger sets the logger to use
func OptLogger(l *zap.SugaredLogger) PlugOption {
	return func(c *plugConfig) {
		c.logger = l
	}
}

// OptHabxEnv sets the habx env to use
func OptHabxEnv(habxEnv string) PlugOption {
	return func(c *plugConfig) {
		c.habxEnv = habxEnv
	}
}
