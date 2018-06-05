package http

import (
	// stdlib
	"context"
	"encoding/json"
	"errors"
	"net/http"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
)

// decodeLoginResponse decodes the incoming HTTP payload to the Go kit payload
func decodeLoginResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.LoginResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeEventCreateResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventCreateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.EventCreateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeEventGetResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventGetResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.EventGetResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeEventUpdateResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventUpdateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.EventUpdateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeEventDeleteResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventDeleteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.EventDeleteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeEventListResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventListResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.EventListResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeUnlockDeviceResponse decodes the incoming HTTP payload to the Go kit payload
func decodeUnlockDeviceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.UnlockDeviceResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeGenerateQRResponse decodes the incoming HTTP payload to the Go kit payload
func decodeGenerateQRResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var resp transport.GenerateQRResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
