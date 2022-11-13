// Package cors provide a default CORS handling code
package cors

import (
	"regexp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/habx/lib-go-utils/generics"
)

// ValidateOrigins validates the origins of a request
func ValidateOrigins(validOrigins []string) func(string) bool {
	reOrigins := generics.Map(validOrigins, regexp.MustCompile)

	return func(origin string) bool {
		for _, r := range reOrigins {
			matched := r.MatchString(origin)
			if matched {
				return true
			}
		}

		return false
	}
}

type corsConfig struct {
	forceTransmission bool
}

// Option is a function that configures the cors middleware
type Option func(c *corsConfig)

// OptForceTransmission forces the transmission of the CORS headers
func OptForceTransmission() Option {
	return func(c *corsConfig) {
		c.forceTransmission = true
	}
}

func createCorsConfig() cors.Config {
	corsConfig := cors.Config{
		// AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},

		// Note: I don't think this configuration makes sense. We most probably don't need to specify all these headers
		AllowHeaders: []string{
			// Standard headers
			"Host", "User-Agent", "Cookie",
			"Origin",
			"Content-Length", "Content-Type",
			"Accept", "Accept-Encoding", "Accept-Language",
			"Cache-Control", "Pragma",
			"Referrer",

			// Security headers
			"Sec-Fetch-Mode", "Sec-Fetch-Site",

			// Apollo headers
			"Apollographql-Client-Name", "Apollographql-Client-Version", "x-apollo-tracing",

			// Habx headers
			"Auth-Token", "X-Request-From", "X-Request-By",
		},

		AllowOriginFunc: ValidateOrigins([]string{
			// Don't forget this: https://regex101.com/
			`http://localhost:.*`,
			`http://puppeteer.*`,
			`https://www\.habx(-dev|-staging)?\.com`,
		}),

		AllowCredentials: true,
	}

	return corsConfig
}

// Handler returns a gin handler for CORS in the habx stack
func Handler(options ...Option) gin.HandlerFunc {
	conf := corsConfig{}

	for _, option := range options {
		option(&conf)
	}

	handler := cors.New(createCorsConfig())

	if conf.forceTransmission {
		previousHandler := handler
		handler = func(c *gin.Context) {
			c.Request.Host = "fake-host-to-force-cors-transmission"
			previousHandler(c)
		}
	}

	return handler
}
