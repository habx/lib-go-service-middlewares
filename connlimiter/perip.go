package connlimiter

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type perIPQueuer struct {
	connectionsLimit        int
	counters                map[string]chan struct{}
	countersMutex           sync.RWMutex
	nbConcurrentConnections int32
}

// QueuePerIP returns a middleware that limits the number of concurrent connections per IP
func QueuePerIP(nb int) gin.HandlerFunc {
	limiter := &perIPQueuer{
		connectionsLimit: nb,
		counters:         make(map[string]chan struct{}),
	}

	return limiter.handler
}

func (l *perIPQueuer) getLimiterRead(ip string) chan struct{} {
	l.countersMutex.RLock()
	defer l.countersMutex.RUnlock()

	return l.counters[ip]
}

func (l *perIPQueuer) getLimiter(ip string) chan struct{} {
	counter := l.getLimiterRead(ip)
	if counter != nil {
		return counter
	}

	// Otherwise, we create it. At this point there can be two different connections
	l.countersMutex.Lock()
	defer l.countersMutex.Unlock()

	// Which means we need to check again if the counter already exists
	counter = l.counters[ip]

	if counter != nil {
		return counter
	}

	// if it's just us and we have reach a max number of IP tracked, we reset the counters
	if len(l.counters) >= nbIPTracked && l.nbConcurrentConnections == 1 {
		l.counters = make(map[string]chan struct{})
	}

	counter = make(chan struct{}, l.connectionsLimit)

	l.counters[ip] = counter

	return counter
}

func (l *perIPQueuer) handler(c *gin.Context) {
	atomic.AddInt32(&l.nbConcurrentConnections, 1)

	ip := c.ClientIP()

	counter := l.getLimiter(ip)

	// before request
	counter <- struct{}{}

	// after request
	defer func() {
		<-counter
		atomic.AddInt32(&l.nbConcurrentConnections, -1)
	}()

	c.Next()
}

const nbIPTracked = 100

// DropPerIP returns a middleware that limits the number of concurrent connections per IP
// by dropping any new connection after that limit
func DropPerIP(n int) gin.HandlerFunc {
	limiter := &perIPDropper{
		connectionsLimit: int32(n),
		counters:         make(map[string]*int32),
	}

	return limiter.handler
}

type perIPDropper struct {
	connectionsLimit        int32
	counters                map[string]*int32
	countersMutex           sync.RWMutex
	nbConcurrentConnections int32
}

func (l *perIPDropper) getLimiterRead(ip string) *int32 {
	l.countersMutex.RLock()
	defer l.countersMutex.RUnlock()

	return l.counters[ip]
}

func (l *perIPDropper) getCounter(ip string) *int32 {
	counter := l.getLimiterRead(ip)
	if counter != nil {
		return counter
	}

	// Otherwise, we create it. At this point there can be two different connections
	l.countersMutex.Lock()
	defer l.countersMutex.Unlock()

	// Which means we need to check again if the counter already exists
	counter = l.counters[ip]

	if counter != nil {
		return counter
	}

	// if it's just us and we have reach a max number of IP tracked, we reset the counters
	if len(l.counters) >= nbIPTracked && l.nbConcurrentConnections == 1 {
		l.counters = make(map[string]*int32)
	}

	counter = new(int32)

	l.counters[ip] = counter

	return counter
}

func (l *perIPDropper) handler(c *gin.Context) {
	ip := c.ClientIP()

	atomic.AddInt32(&l.nbConcurrentConnections, 1)

	counter := l.getCounter(ip)

	// before request
	atomic.AddInt32(counter, 1)

	// after request
	defer atomic.AddInt32(counter, -1)
	defer atomic.AddInt32(&l.nbConcurrentConnections, -1)

	if *counter > l.connectionsLimit {
		c.AbortWithStatus(http.StatusTooManyRequests)

		return
	}

	c.Next()
}
