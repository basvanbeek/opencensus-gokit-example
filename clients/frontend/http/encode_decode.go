package http

import (
	// stdlib
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	// project

	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
)

// decodeLoginResponse decodes the incoming HTTP payload to the Go kit payload
func decodeLoginResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.LoginResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	case http.StatusBadRequest:
		body := decodeErrorResponse(r)
		switch body {
		case frontend.ErrorUserPassRequired:
			resp.Err = frontend.ErrUserPassRequired
		case frontend.ErrorUserPassUnknown:
			resp.Err = frontend.ErrUserPassUnknown
		default:
			return nil, errors.New(body)
		}
		return resp, nil
	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

// decodeEventCreateResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventCreateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.EventCreateResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	case http.StatusConflict:
		body := decodeErrorResponse(r)
		switch body {
		case frontend.ErrorEventExists:
			resp.Err = frontend.ErrEventExists
		default:
			return nil, errors.New(body)
		}
		return resp, nil
	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

// decodeEventGetResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventGetResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.EventGetResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

// decodeEventUpdateResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventUpdateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.EventUpdateResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

// decodeEventDeleteResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventDeleteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.EventDeleteResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

// decodeEventListResponse decodes the incoming HTTP payload to the Go kit payload
func decodeEventListResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.EventListResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

// decodeUnlockDeviceResponse decodes the incoming HTTP payload to the Go kit payload
func decodeUnlockDeviceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.UnlockDeviceResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

// decodeGenerateQRResponse decodes the incoming HTTP payload to the Go kit payload
func decodeGenerateQRResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp transport.GenerateQRResponse

	switch r.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, err
		}
		return resp, nil
	case http.StatusBadRequest:
		body := decodeErrorResponse(r)
		switch body {
		case frontend.ErrorInvalidQRParams:
			resp.Err = frontend.ErrInvalidQRParams
		default:
			return nil, errors.New(body)
		}
		return resp, nil
	case http.StatusServiceUnavailable:
		body := decodeErrorResponse(r)
		switch body {
		case frontend.ErrorQRGenerate:
			return nil, frontend.ErrQRGenerate
		default:
			return nil, errors.New(body)
		}

	default:
		return nil, errors.New(decodeErrorResponse(r))
	}
}

func decodeErrorResponse(r *http.Response) string {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err.Error()
	}

	return string(bytes.Trim(b, "\r\n\t "))
}
