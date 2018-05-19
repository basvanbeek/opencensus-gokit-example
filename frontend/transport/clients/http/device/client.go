package device

import (
	// stdlib
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	// external
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/satori/go.uuid"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/device"
	"github.com/basvanbeek/opencensus-gokit-example/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/device/transport/http/routes"
)

// client http transport to device service.
type client struct {
	endpoints transport.Endpoints
	routes    routes.Endpoints
	logger    log.Logger
}

// New returns a new Device client using http transport
func New(instancer sd.Instancer, logger log.Logger) device.Service {
	options := []kithttp.ClientOption{}

	makeEndpoint := func(
		method string, enc kithttp.EncodeRequestFunc, dec kithttp.DecodeResponseFunc,
	) endpoint.Endpoint {

		// configure circuit breaker
		cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "Device/http/" + method,
			MaxRequests: 5,
			Interval:    10 * time.Second,
			Timeout:     10 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 5
			},
		})

		// configure rate limiter
		rl := rate.NewLimiter(rate.Every(time.Second), 100)

		// our method sd.Factory is called when a new QR service is discovered.
		factory := func(instance string) (endpoint.Endpoint, io.Closer, error) {
			// set-up our go kit client endpoint

			baseURL, err := url.Parse(instance)
			if err != nil {
				return nil, nil, err
			}

			var e endpoint.Endpoint
			e = kithttp.NewClient(method, baseURL, enc, dec, options...).Endpoint()
			// (conn, "pb.QR", method, enc, dec, reply, options...).Endpoint()
			e = ratelimit.NewErroringLimiter(rl)(e)
			e = circuitbreaker.Gobreaker(cb)(e)

			return e, nil, nil
		}

		// endpointer manages list of available endpoints servicing our method
		endpointer := sd.NewEndpointer(instancer, factory, logger)

		// balancer can do a round robin pick from the endpointer list
		balancer := lb.NewRoundRobin(endpointer)

		// retry uses balancer for executing a method call with retry and timeout
		// logic so client consumer does not have to think about it.
		retry := lb.Retry(3, 5*time.Second, balancer)
		return retry
	}

	c := &client{
		logger: logger,
	}

	c.endpoints = transport.Endpoints{
		Unlock: makeEndpoint("Unlock", c.encodeUnlockRequest, decodeUnlockResponse),
	}

	return c
}

func (c client) Unlock(ctx context.Context, eventID, deviceID uuid.UUID, code string) (*device.Session, error) {
	res, err := c.endpoints.Unlock(ctx, transport.UnlockRequest{
		EventID:  eventID,
		DeviceID: deviceID,
		Code:     code,
	})
	if err != nil {
		return nil, err
	}
	response := res.(transport.UnlockResponse)
	return &device.Session{
		EventCaption:  response.EventCaption,
		DeviceCaption: response.DeviceCaption,
	}, nil
}

// encodeUnlockRequest encodes the outgoing go kit payload to the http payload
func (c client) encodeUnlockRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	req := request.(transport.UnlockRequest)
	r.URL, err = c.routes.Unlock.URL(
		"event_id", req.EventID.String(),
		"device_id", req.DeviceID.String(),
		"code", req.Code,
	)
	return
}

// decodeUnlockResponse decodes the incoming grpc payload to go kit payload
func decodeUnlockResponse(_ context.Context, response *http.Response) (interface{}, error) {
	var res transport.UnlockResponse
	dec := json.NewDecoder(response.Body)
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
