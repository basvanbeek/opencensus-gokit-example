package http

import (
	// stdlib
	"context"
	"encoding/json"
	"net/http"

	// external
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	// project

	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/http/routes"
)

// NewHTTPHandler wires our Go kit endpoints to the HTTP transport.
func NewHTTPHandler(svcEndpoints transport.Endpoints) http.Handler {
	// set-up router and initialize http endpoints
	var (
		router        = mux.NewRouter()
		httpEndpoints = routes.InitEndpoints(router)
	)

	// wire our Go kit handlers to the http endpoints
	httpEndpoints.Unlock.Handler(httptransport.NewServer(
		svcEndpoints.Unlock, decodeUnlockRequest, encodeUnlockResponse,
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
