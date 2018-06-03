package frontend

import (
	// stdlib
	"context"
	"time"

	// external
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/sony/gobreaker"
	"go.opencensus.io/trace"
	"golang.org/x/time/rate"

	// project
	fehttp "github.com/basvanbeek/opencensus-gokit-example/services/cli/transport/clients/frontend/http"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/http/routes"
	"github.com/basvanbeek/opencensus-gokit-example/shared/oc"
)

type client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

// NewHTTP returns a new device client using the HTTP transport
func NewHTTP(instancer sd.Instancer, logger log.Logger) frontend.Service {
	// initialize our codec context
	codec := fehttp.Codec{Route: routes.InitEndpoints(mux.NewRouter())}

	// set-up our http transport options
	options := []kithttp.ClientOption{
		kitoc.HTTPClientTrace(), // add OpenCensus Client tracing
	}

	// configure rate limiter
	rl := rate.NewLimiter(rate.Every(time.Second), 1000)

	// configure circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "CLI/http",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})

	// initialize our generic service endpoint middleware
	mw := endpoint.Chain(
		ratelimit.NewErroringLimiter(rl),
		circuitbreaker.Gobreaker(cb),
	)

	return &client{
		endpoints: transport.Endpoints{
			Login: factory(instancer, "Login")(fehttp.NewFactory(
				codec.Login,
				oc.ChainMW("Login", mw), // chained custom method middleware
				options...,
			)),
			EventCreate: factory(instancer, "EventCreate")(fehttp.NewFactory(
				codec.EventCreate,
				oc.ChainMW("EventCreate", mw), // chained custom method middleware
				options...,
			)),
			EventGet: factory(instancer, "EventGet")(fehttp.NewFactory(
				codec.EventGet,
				oc.ChainMW("EventGet", mw), // chained custom method middleware
				options...,
			)),
			EventUpdate: factory(instancer, "EventUpdate")(fehttp.NewFactory(
				codec.EventUpdate,
				oc.ChainMW("EventUpdate", mw), // chained custom method middleware
				options...,
			)),
			EventDelete: factory(instancer, "EventDelete")(fehttp.NewFactory(
				codec.EventDelete,
				oc.ChainMW("EventDelete", mw), // chained custom method middleware
				options...,
			)),
			EventList: factory(instancer, "EventList")(fehttp.NewFactory(
				codec.EventList,
				oc.ChainMW("EventList", mw), // chained custom method middleware
				options...,
			)),
			UnlockDevice: factory(instancer, "UnlockDevice")(fehttp.NewFactory(
				codec.UnlockDevice,
				oc.ChainMW("UnlockDevice", mw), // chained custom method middleware
				options...,
			)),
			GenerateQR: factory(instancer, "GenerateQR")(fehttp.NewFactory(
				codec.GenerateQR,
				oc.ChainMW("GenerateQR", mw), // chained custom method middleware
				options...,
			)),
		},
		logger: logger,
	}
}

func (c client) Login(ctx context.Context, user, pass string) (*frontend.Login, error) {
	response, err := c.endpoints.Login(
		ctx,
		transport.LoginRequest{
			User: user,
			Pass: pass,
		},
	)
	if err != nil {
		return nil, err
	}

	res := response.(transport.LoginResponse)
	if res.Failed() != nil {
		return nil, res.Failed()
	}

	return &frontend.Login{
		ID:         res.ID,
		Name:       res.Name,
		TenantID:   res.TenantID,
		TenantName: res.TenantName,
	}, nil
}

func (c client) EventCreate(ctx context.Context, tenantID uuid.UUID, event frontend.Event) (*uuid.UUID, error) {
	response, err := c.endpoints.EventCreate(
		ctx,
		transport.EventCreateRequest{
			TenantID: tenantID,
			Event:    event,
		},
	)
	if err != nil {
		return nil, err
	}
	res := response.(transport.EventCreateResponse)
	if res.Failed() != nil {
		return nil, err
	}

	return res.EventID, nil
}

func (c client) EventGet(ctx context.Context, tenantID, eventID uuid.UUID) (*frontend.Event, error) {
	response, err := c.endpoints.EventGet(
		ctx,
		transport.EventGetRequest{
			TenantID: tenantID,
			EventID:  eventID,
		},
	)
	if err != nil {
		return nil, err
	}
	res := response.(transport.EventGetResponse)
	if res.Failed() != nil {
		return nil, res.Failed()
	}
	return res.Event, nil
}

func (c client) EventUpdate(ctx context.Context, tenantID uuid.UUID, event frontend.Event) error {
	response, err := c.endpoints.EventUpdate(
		ctx,
		transport.EventUpdateRequest{
			TenantID: tenantID,
			Event:    event,
		},
	)
	if err != nil {
		return err
	}

	return response.(transport.EventUpdateResponse).Failed()
}

func (c client) EventDelete(ctx context.Context, tenantID, eventID uuid.UUID) error {
	response, err := c.endpoints.EventDelete(
		ctx,
		transport.EventDeleteRequest{
			TenantID: tenantID,
			EventID:  eventID,
		},
	)
	if err != nil {
		return err
	}

	return response.(transport.EventDeleteResponse).Failed()
}

func (c client) EventList(ctx context.Context, tenantID uuid.UUID) ([]*frontend.Event, error) {
	response, err := c.endpoints.EventList(
		ctx,
		transport.EventListRequest{
			TenantID: tenantID,
		},
	)
	if err != nil {
		return nil, err
	}

	res := response.(transport.EventListResponse)
	if res.Failed() != nil {
		return nil, res.Failed()
	}

	return res.Events, nil
}

func (c client) UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (*frontend.Session, error) {
	response, err := c.endpoints.UnlockDevice(
		ctx,
		transport.UnlockDeviceRequest{
			EventID:    eventID,
			DeviceID:   deviceID,
			UnlockCode: unlockCode,
		},
	)
	if err != nil {
		return nil, err
	}

	res := response.(transport.UnlockDeviceResponse)
	if res.Failed() != nil {
		return nil, err
	}

	return res.Session, nil
}

func (c client) GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error) {
	response, err := c.endpoints.GenerateQR(
		ctx,
		transport.GenerateQRRequest{
			EventID:    eventID,
			DeviceID:   deviceID,
			UnlockCode: unlockCode,
		},
	)
	if err != nil {
		return nil, err
	}

	res := response.(transport.GenerateQRResponse)
	if res.Failed() != nil {
		return nil, res.Failed()
	}

	return res.QR, nil
}

// factory creates a service discovery driven Go kit client endpoint
func factory(instancer sd.Instancer, opName string) func(sd.Factory) endpoint.Endpoint {
	return func(factory sd.Factory) endpoint.Endpoint {
		// endpointer manages list of available endpoints servicing our method
		endpointer := sd.NewEndpointer(instancer, factory, log.NewNopLogger())

		// balancer can do a random pick from the endpointer list
		balancer := lb.NewRandom(endpointer, time.Now().UnixNano())

		// retry uses balancer for executing a method call with retry and
		// timeout logic so client consumer does not have to think about it.
		var (
			count    = 3
			duration = 5 * time.Second
		)
		endpoint := lb.Retry(count, duration, balancer)

		//endpoint = requestid.Endpoint(true)(endpoint)

		return kitoc.TraceEndpoint(
			"kit/retry "+opName,
			kitoc.WithEndpointAttributes(
				trace.StringAttribute("kit.balancer.type", "random"),
				trace.StringAttribute("kit.retry.timeout", duration.String()),
				trace.Int64Attribute("kit.retry.count", int64(count)),
			),
		)(endpoint)

	}
}
