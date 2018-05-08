package svchttp

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/basvanbeek/opencensus-gokit-example/frontend/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/frontend/transport/svchttp/routes"
)

func NewHTTPHandler(svcEndpoints implementation.Endpoints) http.Handler {
	// set-up router and initialize http endpoints
	var (
		router        = mux.NewRouter()
		httpEndpoints = routes.InitEndpoints(router)
	)

	// wire our Go kit handlers to the http endpoints
	httpEndpoints.UnlockDevice.Handler(httptransport.NewServer(
		svcEndpoints.UnlockDevice, decodeUnlockDeviceRequest, encodeUnlockDeviceResponse,
	))

	httpEndpoints.GenerateQR.Handler(httptransport.NewServer(
		svcEndpoints.GenerateQR, decodeGenerateQRRequest, encodeGenerateQRResponse,
	))

	// return our router as http handler
	return router
}

// decode / encode functions for converting between http transport payloads and
// Go kit request_response payloads.

func decodeUnlockDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.UnlockDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeUnlockDeviceResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeGenerateQRRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.GenerateQRRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeGenerateQRResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(implementation.GenerateQRResponse)
	if res.Failed() != nil {
		// TODO: add logic ex. auth
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		b, err := json.Marshal(res.Failed().Error())
		if err != nil {
			return err
		}
		w.Write(b)
		return nil
	}

	return json.NewEncoder(w).Encode(res.QR)
}
