package frontend

import (
	// stdlib
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/http/routes"
)

// CodecFunc holds codec details for a kit http transport client
type CodecFunc func() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc)

// Codec holds context for our transport encoders & decoders
type Codec struct {
	Route routes.Endpoints
}

// NewFactory returns a new endpoint factory using the HTTP transport for our
// device client.
func NewFactory(
	codecFunc CodecFunc, mw endpoint.Middleware, options ...kithttp.ClientOption,
) sd.Factory {
	// retrieve our codecs
	method, enc, dec := codecFunc()

	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		baseURL, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}

		// set-up our go kit client endpoint
		endpoint := kithttp.NewClient(
			method, baseURL, enc, dec, options...,
		).Endpoint()

		return mw(endpoint), nil, nil
	}
}

// Login codec
func (c Codec) Login() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.Login, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.LoginResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "Login", encRequest, decResponse
}

// EventCreate codec
func (c Codec) EventCreate() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.EventCreate, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.EventCreateResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "EventCreate", encRequest, decResponse
}

// EventGet codec
func (c Codec) EventGet() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.EventGet, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.EventGetResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "EventGet", encRequest, decResponse
}

// EventUpdate codec
func (c Codec) EventUpdate() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.EventUpdate, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.EventUpdateResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "EventUpdate", encRequest, decResponse
}

// EventDelete codec
func (c Codec) EventDelete() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.EventDelete, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.EventDeleteResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "EventDelete", encRequest, decResponse
}

// EventList codec
func (c Codec) EventList() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.EventList, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.EventListResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "EventList", encRequest, decResponse
}

// UnlockDevice codec
func (c Codec) UnlockDevice() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.UnlockDevice, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.UnlockDeviceResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "UnlockDevice", encRequest, decResponse
}

// GenerateQR codec
func (c Codec) GenerateQR() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		return genericEncoder(c.Route.GenerateQR, r, request)
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		var resp transport.GenerateQRResponse
		err := json.NewDecoder(r.Body).Decode(&resp)
		return resp, err
	}

	return "GenerateQR", encRequest, decResponse
}

func genericEncoder(route *mux.Route, r *http.Request, request interface{}) error {
	var (
		err error
		buf bytes.Buffer
	)

	if r.URL, err = route.Host(r.URL.Host).URL(); err != nil {
		return err
	}
	if methods, err := route.GetMethods(); err == nil {
		r.Method = methods[0]
	}

	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
