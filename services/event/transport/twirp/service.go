package twirp

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/kevinburke/go.uuid"
	"github.com/twitchtv/twirp"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/transport/pb"
)

type server struct {
	svc    event.Service
	logger log.Logger
}

// NewService returns a new Service backed by Twirp transport.
func NewService(svc event.Service, logger log.Logger) pb.Event {
	return &server{
		svc:    svc,
		logger: logger,
	}
}

func (s *server) Create(ctx context.Context, r *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := s.svc.Create(
		ctx,
		uuid.FromBytesOrNil(r.TenantId),
		event.Event{
			ID:   uuid.FromBytesOrNil(r.Event.Id),
			Name: r.Event.Name,
		},
	)

	switch err {
	case nil:
		return &pb.CreateResponse{Id: id.Bytes()}, nil
	case event.ErrEventExists:
		return nil, twirp.NewError(twirp.AlreadyExists, err.Error())
	default:
		return nil, twirp.InternalErrorWith(err)
	}
}

func (s *server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	evt, err := s.svc.Get(
		ctx,
		uuid.FromBytesOrNil(r.TenantId),
		uuid.FromBytesOrNil(r.Id),
	)

	switch err {
	case nil:
		return &pb.GetResponse{
			Event: &pb.EventObj{Id: evt.ID.Bytes(), Name: evt.Name},
		}, nil
	case event.ErrNotFound:
		return nil, twirp.NotFoundError(err.Error())
	default:
		return nil, twirp.InternalErrorWith(err)
	}
}

func (s *server) Update(ctx context.Context, r *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	err := s.svc.Update(
		ctx,
		uuid.FromBytesOrNil(r.TenantId),
		event.Event{
			ID:   uuid.FromBytesOrNil(r.Event.Id),
			Name: r.Event.Name,
		},
	)

	switch err {
	case nil:
		return &pb.UpdateResponse{}, nil
	case event.ErrNotFound:
		return nil, twirp.NotFoundError(err.Error())
	case event.ErrEventExists:
		return nil, twirp.NewError(twirp.AlreadyExists, err.Error())
	default:
		return nil, twirp.InternalErrorWith(err)
	}
}

func (s *server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.svc.Delete(
		ctx,
		uuid.FromBytesOrNil(r.TenantId),
		uuid.FromBytesOrNil(r.Id),
	)

	switch err {
	case nil:
		return &pb.DeleteResponse{}, nil
	case event.ErrNotFound:
		return nil, twirp.NotFoundError(err.Error())
	default:
		return nil, twirp.InternalErrorWith(err)
	}
}

func (s *server) List(ctx context.Context, r *pb.ListRequest) (*pb.ListResponse, error) {
	events, err := s.svc.List(ctx, uuid.FromBytesOrNil(r.TenantId))
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}
	pbEvents := make([]*pb.EventObj, 0, len(events))
	for _, event := range events {
		pbEvent := &pb.EventObj{
			Id:   event.ID.Bytes(),
			Name: event.Name,
		}
		pbEvents = append(pbEvents, pbEvent)
	}
	return &pb.ListResponse{Events: pbEvents}, nil
}
