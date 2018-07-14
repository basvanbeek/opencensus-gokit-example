package errormw

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// UnwrapError middleware to move a business error found in the response
// parameter into the error parameter. This function should be called after
// circuit breaking and/or retry logic so a business error will not trigger
// those middlewares.
// It's use is primarily as outer shell middleware at the client end so a client
// method implementation only needs to test for error instead of both error and
// response payload business error.
func UnwrapError(logger log.Logger) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if err != nil {
				logger.Log("err", err, "type", "transport")
				return nil, err
			}
			if f, ok := response.(endpoint.Failer); ok {
				if f.Failed() != nil {
					logger.Log("err", f.Failed(), "type", "business")
					return nil, f.Failed()
				}
			}
			return response, nil
		}
	}
}
