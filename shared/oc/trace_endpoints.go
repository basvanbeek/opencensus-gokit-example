package oc

import (
	// stdlib
	"time"

	// external
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	"go.opencensus.io/trace"
)

// BalancerType used in Retry logic
type BalancerType string

// BalancerTypes
const (
	Random     BalancerType = "random"
	RoundRobin BalancerType = "round robin"
)

// ClientEndpoint adds our Endpoint Tracing middleware to the existing client
// side endpoint.
func ClientEndpoint(operationName string) endpoint.Middleware {
	return kitoc.TraceEndpoint(
		"gokit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("gokit.endpoint.type", "client"),
		),
	)
}

// ServerEndpoint adds our Endpoint Tracing middleware to the existing server
// side endpoint.
func ServerEndpoint(operationName string) endpoint.Middleware {
	return kitoc.TraceEndpoint(
		"gokit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("gokit.endpoint.type", "server"),
		),
	)
}

// RetryEndpoint wraps a Go kit lb.Retry endpoint with an annotated span.
func RetryEndpoint(
	operationName string, balancer BalancerType, max int, timeout time.Duration,
) endpoint.Middleware {
	return kitoc.TraceEndpoint("gokit/retry "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("gokit.balancer.type", string(balancer)),
			trace.StringAttribute("gokit.retry.timeout", timeout.String()),
			trace.Int64Attribute("gokit.retry.max_count", int64(max)),
		),
	)
}
