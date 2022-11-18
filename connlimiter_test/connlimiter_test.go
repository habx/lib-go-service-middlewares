package conlimiter_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	tgin "github.com/habx/lib-go-tests/http/gin"

	"github.com/habx/lib-go-service-middlewares/connlimiter"
)

type countersValues struct {
	Current int
	Max     int
}

type countersVerifier struct {
	counters map[string]*countersValues
	mu       sync.Mutex
}

func newCounterVerifier() *countersVerifier {
	return &countersVerifier{counters: make(map[string]*countersValues)}
}

func (c *countersVerifier) getCounter(name string) *countersValues {
	counter := c.counters[name]

	if counter == nil {
		counter = &countersValues{}
		c.counters[name] = counter
	}

	return counter
}

func (c *countersVerifier) entering(ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	counter := c.getCounter(ip)

	counter.Current++

	if counter.Current > c.counters[ip].Max {
		counter.Max = counter.Current
	}
}

func (c *countersVerifier) leaving(ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	counter := c.getCounter(ip)

	counter.Current--
}

func (c *countersVerifier) Handler(ctx *gin.Context) {
	ip := ctx.ClientIP()

	c.entering(ip)
	defer c.leaving(ip)

	ctx.Next()
}

func (c *countersVerifier) Max() int {
	var max int
	for _, counter := range c.counters {
		if counter.Max > max {
			max = counter.Max
		}
	}

	return max
}

func sendRequests(t *testing.T, url string, nbClients int, nbIps int, nbQueries int) *sync.WaitGroup {
	a := assert.New(t)
	wg := sync.WaitGroup{}

	for i := 0; i < nbClients; i++ {
		ip := fmt.Sprintf("10.0.0.%d", i%nbIps+1)

		wg.Add(1)

		go func() {
			for i := 0; i < nbQueries; i++ {
				req, err := http.NewRequest(http.MethodGet, url, nil)
				a.NoError(err)
				req.Header.Set("X-Forwarded-For", ip)
				resp, err := http.DefaultClient.Do(req)
				a.NoError(err)
				a.Equal(http.StatusOK, resp.StatusCode)

				defer resp.Body.Close()
			}
			wg.Done()
		}()
	}

	return &wg
}

func TestGlobal(t *testing.T) {
	a := assert.New(t)

	srv, eng := tgin.GetServerWithGin(t)
	// eng.Use(conlimiter.GlobalLimit(10))

	eng.Use(connlimiter.Global(10))

	counters := newCounterVerifier()
	eng.Use(counters.Handler)

	eng.GET("/test", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 10)
	})

	wg := sendRequests(t, srv.URL("/test"), 100, 1, 20)

	wg.Wait()

	a.Equal(10, counters.Max())
}

func TestPerIp(t *testing.T) {
	a := assert.New(t)

	srv, eng := tgin.GetServerWithGin(t)
	eng.Use(connlimiter.PerIP(3))

	counters := newCounterVerifier()
	eng.Use(counters.Handler)

	eng.GET("/test", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 10)
	})

	wg := sendRequests(t, srv.URL("/test"), 100, 10, 20)

	wg.Wait()

	a.Equal(3, counters.Max())
}
