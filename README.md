# lib-go-service-middlewares

This library aims to be on par with [lib-service-middlewares](https://github.com/habx/lib-service-middlewares) but for Go.

## /mgmt API

This API brings the following routes:
- `/mgmt/version`
- `/mgmt/health`
- `/mgmt/memstats`
- `/mgmt/crash`
- `/mgmt/pprof/`

It can be used with the following options:
- `OptHealthManager(m *uhealth.Manager)`
- `OptACLManager(m *acl.Manager)`
- `OptBuildInfo(bi *tbuild.Info)`
- `OptToken(t tauth.Token)`
- `OptLogger(l *zap.SugaredLogger)`
- `OptHabxEnv(habxEnv string)`

None of them are mandatory but it's best to set them.

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/mgmt"
)

func main() {
    eng := gin.New()

    mgmt.Plug(eng)
}
```

## Sentry
Add sentry handling to the response.
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/sentry"
)

func main() {
    eng := gin.New()
    eng.Use(gin.Recovery()) // This is mandatory to avoid crash...
    eng.use(sentry.Handler()) // ...because the sentry handler repanics
}
```

## CORS handling
Add CORS headers to the response.

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/cors"
)

func main() {
    eng := gin.New()
    // CORS check:
    eng.Use(cors.Handler())

    // CORS cheeck with forced transmission:
    eng.Use(cors.Handler(cors.cors.Handler()()))

    // [...]
}
```

## Request ID
Adds a request ID to the request context and a `X-Request-Id` header to the response.

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/requestid"
)

func main() {
    eng := gin.New()
    eng.Use(requestid.Handler())
    // [...]
}
```

## Connlimiter
Limits the number of connection that can be opened at the same time.

### Queue connections globally
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/connlimiter"
)

func main() {
    eng := gin.New()
    eng.Use(connlimiter.QueueGlobal(10))
}
```

### Drop connections globally
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/connlimiter"
)

func main() {
    eng := gin.New()
    eng.Use(connlimiter.DropGlobal(10))
}
```

### Queue connections per IP
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/connlimiter"
)

func main() {
    eng := gin.New()
    eng.Use(connlimiter.QueuePerIP(10))
}
```

### Drop connetions per IP
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/habx/lib-go-service-middlewares/connlimiter"
)

func main() {
    eng := gin.New()
    eng.Use(connlimiter.DropPerIP(10))
}
```

## Query param to header
```go
import (
    "github.com/habx/lib-go-service-middlewares/querytoheader"
)

func main() {
    eng := gin.New()
    eng.Use(querytoheader.Handler(map[string]string{"param": "header"}))
     // [...]
}
```
