package twirp

import (
	// stdlib
	"context"
	"errors"

	// external
	"github.com/go-kit/kit/log"
	"github.com/kevinburke/go.uuid"
	"github.com/twitchtv/twirp"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/transport/pb"
	"github.com/basvanbeek/opencensus-gokit-example/shared/sd"
)

type client struct {
	instancer func() pb.Event
	logger    log.Logger
}

func (c client) Create(
	ctx context.Context, tenantID uuid.UUID, evt event.Event,
) (*uuid.UUID, error) {
	ci := c.instancer()
	if ci == nil {
		return nil, sd.ErrNoClients
	}

	res, err := ci.Create(ctx, &pb.CreateRequest{
		TenantId: tenantID.Bytes(),
		Event: &pb.EventObj{
			Id:   evt.ID.Bytes(),
			Name: evt.Name,
		},
	})

	if err != nil {
		if twErr, ok := err.(twirp.Error); ok {
			switch twErr.Msg() {
			case event.ErrorService:
				return nil, event.ErrService
			case event.ErrorUnauthorized:
				return nil, event.ErrUnauthorized
			case event.ErrorEventExists:
				return nil, event.ErrEventExists
			default:
				return nil, errors.New(twErr.Msg())
			}
		}
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

	if twErr, ok := err.(twirp.Error); ok {
		switch twErr.Error() {
		case event.ErrorService:
			return nil, event.ErrService
		case event.ErrorUnauthorized:
			return nil, event.ErrUnauthorized
		case event.ErrorNotFound:
			return nil, event.ErrNotFound
		}
	}

	return &event.Event{
		ID:   uuid.FromBytesOrNil(res.Event.Id),
		Name: res.Event.Name,
	}, nil
}

func (c client) Update(
	ctx context.Context, tenantID uuid.UUID, evt event.Event,
) error {
	ci := c.instancer()
	if ci == nil {
		return sd.ErrNoClients
	}

	_, err := ci.Update(ctx, &pb.UpdateRequest{
		TenantId: tenantID.Bytes(),
		Event: &pb.EventObj{
			Id:   evt.ID.Bytes(),
			Name: evt.Name,
		},
	})

	if err != nil {
		if twErr, ok := err.(twirp.Error); ok {
			switch twErr.Error() {
			case event.ErrorService:
				return event.ErrService
			case event.ErrorUnauthorized:
				return event.ErrUnauthorized
			case event.ErrorNotFound:
				return event.ErrNotFound
			case event.ErrorEventExists:
				return event.ErrEventExists
			}
		}
	}

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

	if err != nil {
		if twErr, ok := err.(twirp.Error); ok {
			switch twErr.Error() {
			case event.ErrorService:
				return event.ErrService
			case event.ErrorUnauthorized:
				return event.ErrUnauthorized
			case event.ErrorNotFound:
				return event.ErrNotFound
			}
		}
	}

	return err
}

func (c client) List(
	ctx context.Context, tenantID uuid.UUID,
) ([]*event.Event, error) {
	ci := c.instancer()
	if ci == nil {
		return nil, sd.ErrNoClients
	}

	pbListResponse, err := ci.List(ctx, &pb.ListRequest{
		TenantId: tenantID.Bytes(),
	})

	if err != nil {
		if twErr, ok := err.(twirp.Error); ok {
			switch twErr.Error() {
			case event.ErrorService:
				return nil, event.ErrService
			case event.ErrorUnauthorized:
				return nil, event.ErrUnauthorized
			}
		}
	}

	events := make([]*event.Event, 0, len(pbListResponse.Events))
	for _, evt := range pbListResponse.Events {
		events = append(events, &event.Event{
			ID:   uuid.FromBytesOrNil(evt.Id),
			Name: evt.Name,
		})
	}
	return events, nil
}
