package endpoint

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/golang/time/rate"
	"github.com/sony/gobreaker"
	"time"
)

var (
	rateLimiter    = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))
	circuitBreaker = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))
)

func addEndpointMiddleware(rawEndpoint endpoint.Endpoint, mws ...endpoint.Middleware) endpoint.Endpoint {
	ept := rawEndpoint
	for _, middleware := range mws {
		ept = middleware(rawEndpoint)
	}
	return ept
}
