package http

import (
	// stdlib
	"context"
	"encoding/json"
	"net/http"

	// external
	"github.com/go-kit/kit/log"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/http/routes"
)

// NewHTTPHandler wires our Go kit endpoints to the HTTP transport.
func NewHTTPHandler(svcEndpoints transport.Endpoints, logger log.Logger) http.Handler {
	// set-up router and initialize http endpoints
	var (
		router      = mux.NewRouter()
		route       = routes.Initialize(router)
		errorLogger = httptransport.ServerErrorLogger(logger)
		ocTracing   = kitoc.HTTPServerTrace()
	)

	options := []httptransport.ServerOption{errorLogger, ocTracing}

	// wire our Go kit handlers to the http endpoints
	route.Unlock.Handler(httptransport.NewServer(
		svcEndpoints.Unlock, decodeUnlockRequest, encodeUnlockResponse,
		options...,
	))

	// return our router as http handler
	return router
}

// decode / encode functions for converting between http transport payloads and
// Go kit request_response payloads.

func decodeUnlockRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.UnlockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeUnlockResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
