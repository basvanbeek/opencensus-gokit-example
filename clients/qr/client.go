package qr

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"

	//project
	"github.com/basvanbeek/opencensus-gokit-example/clients/qr/grpc"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport"
)

// NewGRPCClient returns a new qr client using the gRPC transport.
func NewGRPCClient(instancer sd.Instancer, logger log.Logger) qr.Service {
	return &client{
		endpoints: grpc.InitEndpoints(instancer, logger),
		logger:    logger,
	}
}

// client grpc transport to QR service.
type client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

// Generate calls the QR Service Generate method.
func (c client) Generate(
	ctx context.Context, data string, recLevel qr.RecoveryLevel, size int,
) ([]byte, error) {
	// we can also validate parameters before sending the request
	if len(data) == 0 {
		return nil, qr.ErrNoContent
	}
	if recLevel < qr.LevelL || recLevel > qr.LevelH {
		return nil, qr.ErrInvalidRecoveryLevel
	}
	if size > 4096 {
		return nil, qr.ErrInvalidSize
	}

	// call our client side go kit endpoint
	res, err := c.endpoints.Generate(
		ctx,
		transport.GenerateRequest{Data: data, Level: recLevel, Size: size},
	)
	if err != nil {
		return nil, err
	}
	response := res.(transport.GenerateResponse)
	return response.QR, response.Err
}
