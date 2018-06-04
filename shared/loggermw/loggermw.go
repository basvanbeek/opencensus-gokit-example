package loggermw

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// LoggerMiddleware logs our endpoint request response payloads...
func LoggerMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log(
					"request", fmt.Sprintf("%+v", request),
					"response", fmt.Sprintf("%+v", response),
					"error", err,
					"took", time.Since(begin),
				)
			}(time.Now())

			response, err = next(ctx, request)
			return
		}
	}
}
