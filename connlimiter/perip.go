package connlimiter

import (
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type perIPLimiter struct {
	connectionsLimit        int
	nbIPsTracked            int
	counters                map[string]chan struct{}
	countersMutex           sync.RWMutex
	nbConcurrentConnections int32
}

// PerIP returns a middleware that limits the number of concurrent connections per IP
func PerIP(nb int) gin.HandlerFunc {
	limiter := &perIPLimiter{
		connectionsLimit: nb,
		nbIPsTracked:     nb,
		counters:         make(map[string]chan struct{}),
	}

	return limiter.handler
}

func (l *perIPLimiter) getLimiterRead(ip string) chan struct{} {
	l.countersMutex.RLock()
	defer l.countersMutex.RUnlock()

	return l.counters[ip]
}

func (l *perIPLimiter) getLimiter(ip string) chan struct{} {
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
	if len(l.counters) >= l.nbIPsTracked && l.nbConcurrentConnections == 1 {
		l.counters = make(map[string]chan struct{})
	}

	counter = make(chan struct{}, l.connectionsLimit)

	l.counters[ip] = counter

	return counter
}

func (l *perIPLimiter) handler(c *gin.Context) {
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
