package implementation

import (
	// stdlib
	"context"
	"strings"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
)

// service implements frontend.Service
type service struct {
	evtClient event.Service
	devClient device.Service
	qrClient  qr.Service
	logger    log.Logger
}

// NewService creates and returns a new Frontend service instance
func NewService(evtClient event.Service, devClient device.Service, qrClient qr.Service, logger log.Logger) frontend.Service {
	return &service{
		evtClient: evtClient,
		devClient: devClient,
		qrClient:  qrClient,
		logger:    logger,
	}
}

func (s *service) Login(ctx context.Context, user, pass string) (*frontend.Login, error) {
	// NOTE: obviously this is hardcoded due to demo purposes
	user = strings.Trim(user, "\r\n\t ")
	pass = strings.Trim(pass, "\r\n\t ")
	if len(user) == 0 || len(pass) == 0 {
		return nil, frontend.ErrUserPassRequired
	}

	switch {
	case user == "john" && pass == "doe":
		return &frontend.Login{
			ID:         uuid.NewV5(uuid.NamespaceOID, "user.id.1"),
			Name:       "John Doe",
			TenantID:   uuid.NewV5(uuid.NamespaceOID, "tenant.id.1"),
			TenantName: "Acme Corp.",
		}, nil
	case user == "jane" && pass == "doe":
		return &frontend.Login{
			ID:         uuid.NewV5(uuid.NamespaceOID, "user.id.2"),
			Name:       "Jane Doe",
			TenantID:   uuid.NewV5(uuid.NamespaceOID, "tenant.id.2"),
			TenantName: "Evil Inc.",
		}, nil
	}

	return nil, frontend.ErrUserPassUnknown
}

func (s *service) EventCreate(ctx context.Context, tenantID uuid.UUID, evt frontend.Event) (*uuid.UUID, error) {
	return s.evtClient.Create(ctx, tenantID, event.Event(evt))
}

func (s *service) EventGet(ctx context.Context, tenantID, id uuid.UUID) (*frontend.Event, error) {
	evt, err := s.evtClient.Get(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return (*frontend.Event)(evt), nil
}

func (s *service) EventUpdate(ctx context.Context, tenantID uuid.UUID, evt frontend.Event) error {
	return s.evtClient.Update(ctx, tenantID, event.Event(evt))
}

func (s *service) EventDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.evtClient.Delete(ctx, tenantID, id)
}

func (s *service) EventList(ctx context.Context, tenantID uuid.UUID) ([]*frontend.Event, error) {
	evts, err := s.evtClient.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	events := make([]*frontend.Event, 0, len(evts))
	for _, e := range evts {
		events = append(events, (*frontend.Event)(e))
	}
	return events, nil
}

// Unlockdevice returns a new session for allowing device to check-in participants.
func (s *service) UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (*frontend.Session, error) {
	var (
		logger = log.With(s.logger, "method", "UnlockDevice")
	)

	if eventID == uuid.Nil {
		level.Warn(logger).Log("err", frontend.ErrRequireEventID)
		return nil, frontend.ErrRequireEventID
	}
	if deviceID == uuid.Nil {
		level.Warn(logger).Log("err", frontend.ErrRequireDeviceID)
		return nil, frontend.ErrRequireDeviceID
	}

	unlockCode = strings.Trim(unlockCode, "\r\n\t ")

	if unlockCode == "" {
		level.Warn(logger).Log("err", frontend.ErrRequireUnlockCode)
		return nil, frontend.ErrRequireUnlockCode
	}

	session, err := s.devClient.Unlock(ctx, eventID, deviceID, unlockCode)
	if err != nil {
		return nil, err
	}

	return &frontend.Session{
		EventID:       eventID,
		EventCaption:  session.EventCaption,
		DeviceID:      deviceID,
		DeviceCaption: session.DeviceCaption,
		Token:         "TOKEN",
	}, nil
}

// Generate returns a new QR code device unlock image based on the provided details.
func (s *service) GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error) {
	var (
		logger = log.With(s.logger, "method", "GenerateQR")
	)

	if eventID == uuid.Nil {
		level.Warn(logger).Log("err", frontend.ErrRequireEventID)
		return nil, frontend.ErrRequireEventID
	}
	if deviceID == uuid.Nil {
		level.Warn(logger).Log("err", frontend.ErrRequireDeviceID)
		return nil, frontend.ErrRequireDeviceID
	}

	unlockCode = strings.Trim(unlockCode, "\r\n\t ")

	if unlockCode == "" {
		level.Warn(logger).Log("err", frontend.ErrRequireUnlockCode)
		return nil, frontend.ErrRequireUnlockCode
	}

	return s.qrClient.Generate(
		ctx, eventID.String()+":"+deviceID.String()+":"+unlockCode, qr.LevelM, 256,
	)
}
