package device

import (
	// stdlib
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/http/routes"
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

// Unlock returns our http codecs
func (c Codec) Unlock() (string, kithttp.EncodeRequestFunc, kithttp.DecodeResponseFunc) {
	// encRequest encodes the outgoing Go kit payload to the HTTP payload
	encRequest := func(_ context.Context, r *http.Request, request interface{}) error {
		var err error
		req := request.(transport.UnlockRequest)
		r.URL, err = c.Route.Unlock.URL(
			"event_id", req.EventID.String(),
			"device_id", req.DeviceID.String(),
			"code", req.Code,
		)
		return err
	}

	// decResponse decodes the incoming HTTP payload to the Go kit payload
	decResponse := func(_ context.Context, response *http.Response) (interface{}, error) {
		var res transport.UnlockResponse
		dec := json.NewDecoder(response.Body)
		if err := dec.Decode(&res); err != nil {
			return nil, err
		}
		return res, nil
	}

	return "Unlock", encRequest, decResponse
}
