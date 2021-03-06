package device

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/kevinburke/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/clients/device/grpc"
	"github.com/basvanbeek/opencensus-gokit-example/clients/device/http"
	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
)

// NewHTTPClient returns a new device client using the HTTP transport.
func NewHTTPClient(instancer sd.Instancer, logger log.Logger) device.Service {
	return &client{
		endpoints: http.InitEndpoints(instancer, logger),
		logger:    logger,
	}
}

// NewGRPCClient returns a new device client using the gRPC transport
func NewGRPCClient(instancer sd.Instancer, logger log.Logger) device.Service {
	return &client{
		endpoints: grpc.InitEndpoints(instancer, logger),
		logger:    logger,
	}

}

type client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

func (c client) Unlock(ctx context.Context, eventID, deviceID uuid.UUID, code string) (*device.Session, error) {
	res, err := c.endpoints.Unlock(ctx, transport.UnlockRequest{
		EventID:  eventID,
		DeviceID: deviceID,
		Code:     code,
	})
	if err != nil {
		// transport logic / unknown error
		return nil, err
	}

	response := res.(transport.UnlockResponse)
	if response.Err != nil {
		// business logic error
		return nil, err
	}

	return &device.Session{
		EventCaption:  response.EventCaption,
		DeviceCaption: response.DeviceCaption,
	}, nil
}
