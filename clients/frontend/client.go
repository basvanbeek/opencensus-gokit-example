package frontend

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/clients/frontend/http"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
)

// NewHTTPClient returns a new frontend client using the HTTP transport.
func NewHTTPClient(instancer sd.Instancer, logger log.Logger) frontend.Service {
	return &client{
		endpoints: httpclient.InitEndpoints(instancer, logger),
		logger:    logger,
	}
}

type client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

func (c *client) Login(ctx context.Context, user, pass string) (*frontend.Login, error) {
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

func (c *client) EventCreate(ctx context.Context, tenantID uuid.UUID, event frontend.Event) (*uuid.UUID, error) {
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

func (c *client) EventGet(ctx context.Context, tenantID, eventID uuid.UUID) (*frontend.Event, error) {
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

func (c *client) EventUpdate(ctx context.Context, tenantID uuid.UUID, event frontend.Event) error {
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

func (c *client) EventDelete(ctx context.Context, tenantID, eventID uuid.UUID) error {
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

func (c *client) EventList(ctx context.Context, tenantID uuid.UUID) ([]*frontend.Event, error) {
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

func (c *client) UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (*frontend.Session, error) {
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

func (c *client) GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error) {
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
