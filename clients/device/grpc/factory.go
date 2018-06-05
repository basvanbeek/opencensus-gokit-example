package device

import (
	// stdlib
	"context"
	"io"

	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/grpc/pb"
	"github.com/basvanbeek/opencensus-gokit-example/shared/grpcconn"
)

// CodecFunc holds codec details for a kit gRPC transport client
type CodecFunc func() (string, kitgrpc.EncodeRequestFunc, kitgrpc.DecodeResponseFunc, interface{})

// NewFactory returns a new endpoint factory using the gRPC transport for our
// device client.
func NewFactory(
	codecFunc CodecFunc, mw endpoint.Middleware, hm grpcconn.HostMapper,
	options ...kitgrpc.ClientOption,
) sd.Factory {
	// retrieve our codecs
	method, enc, dec, reply := codecFunc()

	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		// try to get connection to advertised instance
		conn, closer, err := hm.Get(instance)
		if err != nil {
			return nil, nil, err
		}

		// set-up our go kit client endpoint
		endpoint := kitgrpc.NewClient(
			conn, "pb.Device", method, enc, dec, reply, options...,
		).Endpoint()

		return mw(endpoint), closer, nil
	}

}

// Unlock returns our grpc codecs
func Unlock() (string, kitgrpc.EncodeRequestFunc, kitgrpc.DecodeResponseFunc, interface{}) {
	// encRequest encodes the outgoing go kit payload to the gRPC payload
	encRequest := func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.UnlockRequest)
		return &pb.UnlockRequest{
			EventId: req.EventID.Bytes(),
		}, nil
	}

	// decResponse decodes the incoming gRPC payload to go kit payload
	decResponse := func(_ context.Context, response interface{}) (interface{}, error) {
		res := response.(*pb.UnlockResponse)
		return transport.UnlockResponse{
			DeviceCaption: res.DeviceCaption,
			EventCaption:  res.EventCaption,
		}, nil
	}

	return "Unlock", encRequest, decResponse, pb.UnlockResponse{}
}
