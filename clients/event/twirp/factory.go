package twirp

import (
	// stdlib
	"io"
	"net/http"

	// external
	"github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/transport/pb"
	"github.com/basvanbeek/opencensus-gokit-example/shared/sd"
)

// NewClient returns a new event client using the Twirp transport.
func NewClient(instancer kitsd.Instancer, c *http.Client, logger log.Logger) event.Service {
	return &client{
		instancer: factory(instancer, c, logger),
		logger:    logger,
	}
}

func factory(instancer kitsd.Instancer, client *http.Client, logger log.Logger) func() pb.Event {
	factoryFunc := func(instance string) (interface{}, io.Closer, error) {
		return pb.NewEventProtobufClient(instance, client), nil, nil
	}
	clientInstancer := sd.NewClientInstancer(instancer, factoryFunc, logger)
	balancer := sd.NewRoundRobin(clientInstancer)

	return func() pb.Event {
		client, err := balancer.Client()
		if err != nil {
			logger.Log("err", err)
			return nil
		}
		return client.(pb.Event)
	}
}
