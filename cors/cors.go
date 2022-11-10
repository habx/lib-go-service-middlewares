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

// Handler returns a gin handler for CORS in the habx stack
func Handler() gin.HandlerFunc {
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

	return cors.New(corsConfig)
}
