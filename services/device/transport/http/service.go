package http

import (
	// stdlib
	"context"
	"encoding/json"
	"net/http"

	// external

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/http/routes"
)

// NewService wires our Go kit endpoints to the HTTP transport.
func NewService(
	svcEndpoints transport.Endpoints, options []kithttp.ServerOption,
	logger log.Logger,
) http.Handler {
	// set-up router and initialize http endpoints
	var (
		router       = mux.NewRouter()
		route        = routes.Initialize(router)
		errorLogger  = kithttp.ServerErrorLogger(logger)
		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
	)

	options = append(options, errorLogger, errorEncoder)

	// wire our Go kit handlers to the http endpoints
	route.Unlock.Handler(kithttp.NewServer(
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

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	var code int
	switch err {
	case device.ErrRequireEventID, device.ErrRequireDeviceID, device.ErrRequireUnlockCode:
		code = http.StatusBadRequest
	case device.ErrEventNotFound, device.ErrUnlockNotFound:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), code)
}
