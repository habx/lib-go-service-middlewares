package mgmt_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	acl "github.com/habx/lib-go-acl/v3"
	tcreds "github.com/habx/lib-go-tests/credentials"
	thttp "github.com/habx/lib-go-tests/http"
	buildt "github.com/habx/lib-go-types/build"
	uhealth "github.com/habx/lib-go-utils/health"

	"github.com/habx/lib-go-service-middlewares/mgmt"
	"github.com/habx/lib-go-service-middlewares/mgmt/crash"
	"github.com/habx/lib-go-service-middlewares/mgmt/health"
	"github.com/habx/lib-go-service-middlewares/mgmt/memstats"
	"github.com/habx/lib-go-service-middlewares/mgmt/pprof"
	"github.com/habx/lib-go-service-middlewares/mgmt/version"
)

func TestGlobal(t *testing.T) {
	a := assert.New(t)

	eng := gin.New()
	eng.Use(gin.Recovery()) // For /mgmt/crash
	a.NoError(mgmt.Plug(eng, mgmt.OptHabxEnv("dev")))

	srv := thttp.GetServer(t, thttp.OptHandler(eng))

	t.Run("version", func(t *testing.T) {
		c := srv.GetClient()
		a.Contains(c.GetString("/mgmt/version"), "\"version\":\"")
	})

	t.Run("health", func(t *testing.T) {
		c := srv.GetClient()
		a.Equal(c.GetString("/mgmt/health"), `{"status":"ok"}`)
	})

	token := tcreds.GetUserToken()

	t.Run("crash", func(t *testing.T) {
		c := srv.GetClient(thttp.OptCltToken(token))
		a.Equal(http.StatusInternalServerError, c.GetStatusCode("/mgmt/crash"))
	})

	t.Run("memstats no auth", func(t *testing.T) {
		c := srv.GetClient()
		a.Equal(http.StatusUnauthorized, c.GetStatusCode("/mgmt/memstats"))
	})

	t.Run("memstats", func(t *testing.T) {
		c := srv.GetClient(thttp.OptCltToken(token))
		a.Contains(c.GetString("/mgmt/memstats"), "Alloc")
	})

	t.Run("pprof no auth", func(t *testing.T) {
		c := srv.GetClient()
		a.Equal(http.StatusUnauthorized, c.GetStatusCode("/mgmt/pprof/symbol"))
	})

	t.Run("pprof bad auth", func(t *testing.T) {
		c := srv.GetClient(thttp.OptCltToken(tcreds.GetTestUserToken()))
		a.Equal(http.StatusUnauthorized, c.GetStatusCode("/mgmt/pprof/symbol"))
	})

	t.Run("pprof", func(t *testing.T) {
		c := srv.GetClient(thttp.OptCltToken(token))
		a.Contains(c.GetString("/mgmt/pprof/symbol"), "num_symbols:")
	})

	t.Run("pprof auth by query", func(t *testing.T) {
		c := srv.GetClient()
		a.Contains(c.GetString("/mgmt/pprof/symbol?token="+token.Value), "num_symbols:")
	})
}

func TestGlobalForHealthWithOption(t *testing.T) {
	a := assert.New(t)

	eng := gin.New()
	eng.Use(gin.Recovery()) // For /mgmt/crash
	a.NoError(mgmt.Plug(eng,
		mgmt.OptHabxEnv("dev"),
		mgmt.OptHealthManager(uhealth.NewManager()),
		mgmt.OptLogger(zap.NewNop().Sugar()),
		mgmt.OptToken(tcreds.GetComponentToken())),
	)

	srv := thttp.GetServer(t, thttp.OptHandler(eng))

	t.Run("health", func(t *testing.T) {
		c := srv.GetClient()
		a.Equal(c.GetString("/mgmt/health"), `{"status":"ok"}`)
	})
}

func TestGlobalForVersionWithOption(t *testing.T) {
	a := assert.New(t)

	eng := gin.New()
	eng.Use(gin.Recovery()) // For /mgmt/c
	a.NoError(mgmt.Plug(eng,
		mgmt.OptHabxEnv("dev"),
		mgmt.OptBuildInfo(&buildt.Info{Version: "1.2.3"}),
	))

	srv := thttp.GetServer(t, thttp.OptHandler(eng))

	t.Run("version", func(t *testing.T) {
		c := srv.GetClient()
		a.Equal(`{"version":"1.2.3"}`, c.GetString("/mgmt/version"))
	})
}

func TestGlobalWithACLManager(t *testing.T) {
	a := assert.New(t)

	aclManager, err := acl.NewACLManager(acl.OptHabxEnv("dev"))
	a.NoError(err)

	eng := gin.New()
	eng.Use(gin.Recovery())
	a.NoError(mgmt.Plug(eng,
		mgmt.OptHabxEnv("dev"),
		mgmt.OptACLManager(aclManager),
	))
}

func TestVersion(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	version.Plug(r, &buildt.Info{Version: "1.2.3"})

	s := thttp.GetServer(t, thttp.OptHandler(r))
	c := s.GetClient()
	a.Contains(c.GetString("/mgmt/version"), "1.2.3")
}

func TestHealth(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	h := uhealth.NewManager()

	{
		count := 0
		h.AddCheck("test", func() error {
			count++
			if count > 1 {
				return errors.New("sample error") //nolint:goerr113
			}

			return nil
		})
	}
	health.Plug(r, h)

	s := thttp.GetServer(t, thttp.OptHandler(r))
	c := s.GetClient()
	a.Equal(c.GetString("/mgmt/health"), `{"status":"ok"}`)
	resp := c.Get("/mgmt/health")

	defer resp.Body.Close()

	a.Equal(http.StatusInternalServerError, resp.StatusCode)
	content, err := io.ReadAll(resp.Body)
	a.NoError(err)
	a.Equal(`{"error":"check test failed: sample error","status":"error"}`, string(content))
}

func TestMemStats(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	memstats.Plug(r)

	s := thttp.GetServer(t, thttp.OptHandler(r))
	c := s.GetClient()
	a.Contains(c.GetString("/mgmt/memstats"), "Alloc")
}

func TestPprof(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	pprof.Plug(r.Group("/"))

	s := thttp.GetServer(t, thttp.OptHandler(r))
	c := s.GetClient()
	a.Contains(c.GetString("/mgmt/pprof/symbol"), "num_symbols:")
}

func TestCrash(t *testing.T) {
	a := assert.New(t)

	r := gin.Default()
	r.Use(gin.Recovery())
	crash.Plug(r)

	s := thttp.GetServer(t, thttp.OptHandler(r))
	c := s.GetClient()
	a.Equal(http.StatusInternalServerError, c.GetStatusCode("/mgmt/crash"))
}
