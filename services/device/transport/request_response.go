package transport

import (
	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/satori/go.uuid"
)

var (
	_ endpoint.Failer = UnlockResponse{}
)

// UnlockRequest holds the request parameters for the Unlock method.
type UnlockRequest struct {
	EventID  uuid.UUID
	DeviceID uuid.UUID
	Code     string
}

// UnlockResponse holds the response values for the Unlock method.
type UnlockResponse struct {
	EventCaption  string
	DeviceCaption string
	Err           error
}

// Failed implements Failer
func (r UnlockResponse) Failed() error { return r.Err }
