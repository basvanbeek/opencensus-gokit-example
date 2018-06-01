package http

import (
	// stdlib
	"context"
	"encoding/json"
	"net/http"

	// external
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/http/routes"
)

// NewHTTPHandler wires our Go kit endpoints to the HTTP transport.
func NewHTTPHandler(svcEndpoints transport.Endpoints, logger log.Logger) http.Handler {
	// set-up router and initialize http endpoints
	var (
		router        = mux.NewRouter()
		httpEndpoints = routes.InitEndpoints(router)
		errorLogger   = httptransport.ServerErrorLogger(logger)
	)

	// wire our Go kit handlers to the http endpoints
	httpEndpoints.Login.Handler(httptransport.NewServer(
		svcEndpoints.Login, decodeLoginRequest, encodeLoginResponse,
		errorLogger,
	))

	httpEndpoints.EventCreate.Handler(httptransport.NewServer(
		svcEndpoints.EventCreate, decodeEventCreateRequest, encodeEventCreateResponse,
		errorLogger,
	))

	httpEndpoints.EventGet.Handler(httptransport.NewServer(
		svcEndpoints.EventGet, decodeEventGetRequest, encodeEventGetResponse,
		errorLogger,
	))

	httpEndpoints.EventUpdate.Handler(httptransport.NewServer(
		svcEndpoints.EventUpdate, decodeEventUpdateRequest, encodeEventUpdateResponse,
		errorLogger,
	))

	httpEndpoints.EventDelete.Handler(httptransport.NewServer(
		svcEndpoints.EventDelete, decodeEventDeleteRequest, encodeEventDeleteResponse,
		errorLogger,
	))

	httpEndpoints.EventList.Handler(httptransport.NewServer(
		svcEndpoints.EventList, decodeEventListRequest, encodeEventListResponse,
		errorLogger,
	))

	httpEndpoints.UnlockDevice.Handler(httptransport.NewServer(
		svcEndpoints.UnlockDevice, decodeUnlockDeviceRequest, encodeUnlockDeviceResponse,
		errorLogger,
	))

	httpEndpoints.GenerateQR.Handler(httptransport.NewServer(
		svcEndpoints.GenerateQR, decodeGenerateQRRequest, encodeGenerateQRResponse,
		errorLogger,
	))

	// return our router as http handler
	return router
}

// decode / encode functions for converting between http transport payloads and
// Go kit request_response payloads.

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := response.(transport.Failer).Failed(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeEventCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EventCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventCreateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := response.(transport.Failer).Failed(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}
func decodeEventGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EventGetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventGetResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := response.(transport.Failer).Failed(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}
func decodeEventUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EventUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventUpdateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := response.(transport.Failer).Failed(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}
func decodeEventDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EventDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventDeleteResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := response.(transport.Failer).Failed(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}
func decodeEventListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EventListRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeEventListResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := response.(transport.Failer).Failed(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeUnlockDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.UnlockDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeUnlockDeviceResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := response.(transport.Failer).Failed(); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeGenerateQRRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		err error
		req transport.GenerateQRRequest
	)
	v := mux.Vars(r)
	if req.EventID, err = uuid.FromString(v["event_id"]); err != nil {
		return nil, err
	}

	if req.DeviceID, err = uuid.FromString(v["device_id"]); err != nil {
		return nil, err
	}
	req.UnlockCode = v["code"]
	return req, nil
}

func encodeGenerateQRResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(transport.GenerateQRResponse)
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
