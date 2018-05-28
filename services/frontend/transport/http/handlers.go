package http

import (
	// stdlib
	"context"
	"encoding/json"
	"net/http"

	// external
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/http/routes"
)

// NewHTTPHandler wires our Go kit endpoints to the HTTP transport.
func NewHTTPHandler(svcEndpoints implementation.Endpoints) http.Handler {
	// set-up router and initialize http endpoints
	var (
		router        = mux.NewRouter()
		httpEndpoints = routes.InitEndpoints(router)
	)

	// wire our Go kit handlers to the http endpoints
	httpEndpoints.Login.Handler(httptransport.NewServer(
		svcEndpoints.Login, decodeLoginRequest, encodeLoginResponse,
	))

	httpEndpoints.EventCreate.Handler(httptransport.NewServer(
		svcEndpoints.EventCreate, decodeEventCreateRequest, encodeEventCreateResponse,
	))

	httpEndpoints.EventGet.Handler(httptransport.NewServer(
		svcEndpoints.EventGet, decodeEventGetRequest, encodeEventGetResponse,
	))

	httpEndpoints.EventUpdate.Handler(httptransport.NewServer(
		svcEndpoints.EventUpdate, decodeEventUpdateRequest, encodeEventUpdateResponse,
	))

	httpEndpoints.EventDelete.Handler(httptransport.NewServer(
		svcEndpoints.EventDelete, decodeEventDeleteRequest, encodeEventDeleteResponse,
	))

	httpEndpoints.EventList.Handler(httptransport.NewServer(
		svcEndpoints.EventList, decodeEventListRequest, encodeEventListResponse,
	))

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

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(implementation.LoginResponse)
	if err := res.Failed(); err != nil {
		// create new auth header
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeEventCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.EventCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventCreateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func decodeEventGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.EventGetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventGetResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func decodeEventUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.EventUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventUpdateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func decodeEventDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.EventDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventDeleteResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func decodeEventListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req implementation.EventListRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

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
	v := mux.Vars(r)
	req.EventID = uuid.FromStringOrNil(v["event_id"])
	req.DeviceID = uuid.FromStringOrNil(v["device_id"])
	req.UnlockCode = v["code"]
	return req, nil
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
	w.Header().Set("Content-Type", "image/png")
	w.Write(res.QR)
	return nil
}
