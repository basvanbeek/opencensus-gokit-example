package event

import (
	// stdlib
	"context"
	"io"
	"net/http"

	// external
	"github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/transport/pb"
	"github.com/basvanbeek/opencensus-gokit-example/shared/sd"
)

type client struct {
	instancer func() pb.Event
	logger    log.Logger
}

// NewClient returns a new event client using the Twirp transport.
func NewClient(instancer kitsd.Instancer, logger log.Logger) event.Service {
	return &client{
		instancer: factory(instancer, logger),
		logger:    logger,
	}
}

// NewFactory returns a new client instance factory using the Twirp transport
// for our event client.
func NewFactory() sd.Factory {
	return func(instance string) (interface{}, io.Closer, error) {
		client := pb.NewEventProtobufClient(instance, &http.Client{})
		return client, nil, nil
	}
}

func (c client) Create(
	ctx context.Context, tenantID uuid.UUID, event event.Event,
) (*uuid.UUID, error) {
	ci := c.instancer()
	if ci == nil {
		return nil, sd.ErrNoClients
	}

	res, err := ci.Create(ctx, &pb.CreateRequest{
		TenantId: tenantID.Bytes(),
		Event: &pb.EventObj{
			Id:   event.ID.Bytes(),
			Name: event.Name,
		},
	})
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromBytes(res.Id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (c client) Get(
	ctx context.Context, tenantID, id uuid.UUID,
) (*event.Event, error) {
	ci := c.instancer()
	if ci == nil {
		return nil, sd.ErrNoClients
	}

	res, err := ci.Get(ctx, &pb.GetRequest{
		TenantId: tenantID.Bytes(),
		Id:       id.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	return &event.Event{
		ID:   uuid.FromBytesOrNil(res.Event.Id),
		Name: res.Event.Name,
	}, nil
}

func (c client) Update(
	ctx context.Context, tenantID uuid.UUID, event event.Event,
) error {
	ci := c.instancer()
	if ci == nil {
		return sd.ErrNoClients
	}

	_, err := ci.Update(ctx, &pb.UpdateRequest{
		TenantId: tenantID.Bytes(),
		Event: &pb.EventObj{
			Id:   event.ID.Bytes(),
			Name: event.Name,
		},
	})

	return err
}

func (c client) Delete(
	ctx context.Context, tenantID uuid.UUID, id uuid.UUID,
) error {
	ci := c.instancer()
	if ci == nil {
		return sd.ErrNoClients
	}

	_, err := ci.Delete(ctx, &pb.DeleteRequest{
		TenantId: tenantID.Bytes(),
		Id:       id.Bytes(),
	})

	return err
}

func (c client) List(
	ctx context.Context, tenantID uuid.UUID,
) ([]*event.Event, error) {
	ci := c.instancer()
	if ci == nil {
		return nil, sd.ErrNoClients
	}

	return nil, nil
}

func factory(instancer kitsd.Instancer, logger log.Logger) func() pb.Event {
	return func() pb.Event {
		factory := func(instance string) (interface{}, io.Closer, error) {
			cl := pb.NewEventProtobufClient(instance, &http.Client{})
			return cl, nil, nil
		}

		clientInstancer := sd.NewClientInstancer(instancer, factory, logger)

		balancer := sd.NewRoundRobin(clientInstancer)

		client, err := balancer.Client()
		if err != nil {
			logger.Log("err", err)
			return nil
		}
		return client.(pb.Event)
	}
}
