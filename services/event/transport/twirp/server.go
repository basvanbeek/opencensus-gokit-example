package twirp

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	uuid "github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/transport/pb"
)

type server struct {
	svc    event.Service
	logger log.Logger
}

// NewTwirpServer returns a new Service backed by Twirp transport.
func NewTwirpServer(svc event.Service, logger log.Logger) pb.Event {
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
	if err != nil {
		return nil, err
	}
	return &pb.CreateResponse{Id: id.Bytes()}, nil
}

func (s *server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	event, err := s.svc.Get(
		ctx,
		uuid.FromBytesOrNil(r.TenantId),
		uuid.FromBytesOrNil(r.Id),
	)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Event: &pb.EventObj{
			Id:   event.ID.Bytes(),
			Name: event.Name,
		},
	}, nil
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
	return nil, err
}

func (s *server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.svc.Delete(
		ctx,
		uuid.FromBytesOrNil(r.TenantId),
		uuid.FromBytesOrNil(r.Id),
	)
	return nil, err
}

func (s *server) List(ctx context.Context, r *pb.ListRequest) (*pb.ListResponse, error) {
	events, err := s.svc.List(ctx, uuid.FromBytesOrNil(r.TenantId))
	if err != nil {
		return nil, err
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
