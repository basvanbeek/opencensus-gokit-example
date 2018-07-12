package http

import (
	// stdlib
	"context"
	"encoding/json"
	"net/http"

	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kevinburke/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/http/routes"
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
	route.Login.Handler(kithttp.NewServer(
		svcEndpoints.Login, decodeLoginRequest, encodeLoginResponse,
		options...,
	))

	route.EventCreate.Handler(kithttp.NewServer(
		svcEndpoints.EventCreate, decodeEventCreateRequest, encodeEventCreateResponse,
		options...,
	))

	route.EventGet.Handler(kithttp.NewServer(
		svcEndpoints.EventGet, decodeEventGetRequest, encodeEventGetResponse,
		options...,
	))

	route.EventUpdate.Handler(kithttp.NewServer(
		svcEndpoints.EventUpdate, decodeEventUpdateRequest, encodeEventUpdateResponse,
		options...,
	))

	route.EventDelete.Handler(kithttp.NewServer(
		svcEndpoints.EventDelete, decodeEventDeleteRequest, encodeEventDeleteResponse,
		options...,
	))

	route.EventList.Handler(kithttp.NewServer(
		svcEndpoints.EventList, decodeEventListRequest, encodeEventListResponse,
		options...,
	))

	route.UnlockDevice.Handler(kithttp.NewServer(
		svcEndpoints.UnlockDevice, decodeUnlockDeviceRequest, encodeUnlockDeviceResponse,
		options...,
	))

	route.GenerateQR.Handler(kithttp.NewServer(
		svcEndpoints.GenerateQR, decodeGenerateQRRequest, encodeGenerateQRResponse,
		options...,
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
	if err := response.(endpoint.Failer).Failed(); err != nil {
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
	if err := response.(endpoint.Failer).Failed(); err != nil {
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
	if err := response.(endpoint.Failer).Failed(); err != nil {
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
	if err := response.(endpoint.Failer).Failed(); err != nil {
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
	if err := response.(endpoint.Failer).Failed(); err != nil {
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
	if err := response.(endpoint.Failer).Failed(); err != nil {
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
	if err := response.(endpoint.Failer).Failed(); err != nil {
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

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	var code int
	switch err {
	case frontend.ErrUserPassRequired, frontend.ErrRequireEventID,
		frontend.ErrRequireDeviceID, frontend.ErrRequireUnlockCode:
		code = http.StatusBadRequest
	case frontend.ErrEventNotFound:
		code = http.StatusNotFound
	case frontend.ErrUserPassUnknown, frontend.ErrUnlockNotFound:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), code)
}
