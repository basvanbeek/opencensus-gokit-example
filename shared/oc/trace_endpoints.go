package oc

import (
	// external
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	"go.opencensus.io/trace"
)

// ChainMW adds our Client Endpoint Tracing middleware to the existing
// middleware chain.
// TODO: deprecate in favor of using ClientEndpoint in cleaner factory
func ChainMW(opName string, mw endpoint.Middleware) endpoint.Middleware {
	attrs := kitoc.WithEndpointAttributes(
		trace.StringAttribute("kit.endpoint.type", "client"),
	)
	t := kitoc.TraceEndpoint("kit/endpoint "+opName, attrs)

	return endpoint.Chain(t, mw)
}

// ClientEndpoint adds our Endpoint Tracing middleware to the existing client
// side endpoint.
func ClientEndpoint(operationName string) endpoint.Middleware {
	return kitoc.TraceEndpoint(
		"kit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("kit.endpoint.type", "client"),
		),
	)
}

// ServerEndpoint adds our Endpoint Tracing middleware to the existing server
// side endpoint.
func ServerEndpoint(operationName string) endpoint.Middleware {
	return kitoc.TraceEndpoint(
		"kit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("kit.endpoint.type", "server"),
		),
	)
}
