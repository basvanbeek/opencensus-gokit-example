package transport

import (
	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/kevinburke/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
)

var (
	_ endpoint.Failer = LoginResponse{}
	_ endpoint.Failer = EventCreateResponse{}
	_ endpoint.Failer = EventGetResponse{}
	_ endpoint.Failer = EventUpdateResponse{}
	_ endpoint.Failer = EventDeleteResponse{}
	_ endpoint.Failer = EventListResponse{}
	_ endpoint.Failer = UnlockDeviceResponse{}
	_ endpoint.Failer = GenerateQRResponse{}
)

// LoginRequest holds the request parameters for the Login method.
type LoginRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// LoginResponse holds the response values for the Login method.
type LoginResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	TenantID   uuid.UUID `json:"tenant_id"`
	TenantName string    `json:"tenant_name"`
	Err        error
}

// Failed implements Failer.
func (r LoginResponse) Failed() error { return r.Err }

// EventCreateRequest holds the request parameters for the EventCreate method.
type EventCreateRequest struct {
	TenantID uuid.UUID      `json:"tenant_id"`
	Event    frontend.Event `json:"event"`
}

// EventCreateResponse holds the response values for the EventCreate method.
type EventCreateResponse struct {
	EventID *uuid.UUID `json:"event_id,omitempty"`
	Err     error
}

// Failed implements Failer.
func (r EventCreateResponse) Failed() error { return r.Err }

// EventGetRequest holds the request parameters for the EventGet method.
type EventGetRequest struct {
	TenantID uuid.UUID `json:"tenant_id"`
	EventID  uuid.UUID `json:"event_id"`
}

// EventGetResponse holds the response values for the EventGet method.
type EventGetResponse struct {
	Event *frontend.Event `json:"event,omitempty"`
	Err   error
}

// Failed implements Failer.
func (r EventGetResponse) Failed() error { return r.Err }

// EventUpdateRequest holds the request parameters for the EventUpdate method.
type EventUpdateRequest struct {
	TenantID uuid.UUID      `json:"tenant_id"`
	Event    frontend.Event `json:"event"`
}

// EventUpdateResponse holds the response values for the EventUpdate method.
type EventUpdateResponse struct {
	Err error
}

// Failed implements Failer.
func (r EventUpdateResponse) Failed() error { return r.Err }

// EventDeleteRequest holds the request parameters for the EventDelete method.
type EventDeleteRequest struct {
	TenantID uuid.UUID `json:"tenant_id"`
	EventID  uuid.UUID `json:"event_id"`
}

// EventDeleteResponse holds the response values for the EventDelete method.
type EventDeleteResponse struct {
	Err error
}

// Failed implements Failer.
func (r EventDeleteResponse) Failed() error { return r.Err }

// EventListRequest holds the request parameters for the EventList method.
type EventListRequest struct {
	TenantID uuid.UUID `json:"tenant_id"`
}

// EventListResponse holds the response values for the EventList method.
type EventListResponse struct {
	Events []*frontend.Event `json:"events,omitempty"`
	Err    error
}

// Failed implements Failer.
func (r EventListResponse) Failed() error { return r.Err }

// UnlockDeviceRequest holds the request parameters for the UnlockDevice method.
type UnlockDeviceRequest struct {
	EventID    uuid.UUID `json:"event_id"`
	DeviceID   uuid.UUID `json:"device_id"`
	UnlockCode string    `json:"unlock_code"`
}

// UnlockDeviceResponse holds the response values for the UnlockDevice method.
type UnlockDeviceResponse struct {
	Session *frontend.Session `json:"session,omitempty"`
	Err     error
}

// Failed implements Failer.
func (r UnlockDeviceResponse) Failed() error { return r.Err }

// GenerateQRRequest holds the request parameters for the GenerateQR method.
type GenerateQRRequest struct {
	EventID    uuid.UUID
	DeviceID   uuid.UUID
	UnlockCode string
}

// GenerateQRResponse holds the response values for the GenerateQR method.
type GenerateQRResponse struct {
	QR  []byte
	Err error
}

// Failed implements Failer.
func (r GenerateQRResponse) Failed() error { return r.Err }
