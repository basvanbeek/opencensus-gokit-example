package oc

import (
	// external
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	"go.opencensus.io/trace"
)

// ChainMW adds our Client Endpoint Tracing middleware to the existing
// middleware chain.
func ChainMW(opName string, mw endpoint.Middleware) endpoint.Middleware {
	attrs := kitoc.WithEndpointAttributes(
		trace.StringAttribute("kit.endpoint.type", "client"),
	)
	t := kitoc.TraceEndpoint("kit/endpoint "+opName, attrs)

	return endpoint.Chain(t, mw)
}

// ServerEndpoint adds our Server Endpoint Tracing middleware to the existing
// server endpoint.
func ServerEndpoint(opName string) endpoint.Middleware {
	attrs := kitoc.WithEndpointAttributes(
		trace.StringAttribute("kit.endpoint.type", "server"),
	)
	return kitoc.TraceEndpoint("kit/endpoint "+opName, attrs)
}
