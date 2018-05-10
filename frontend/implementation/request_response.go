package implementation

import (
	// external
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/frontend"
)

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
//
// This is particularly useful if one can abstract the response encoding for
// at least one of the supported service transports (ex. JSON over HTTP).
type Failer interface {
	Failed() error
}

// UnlockDeviceRequest holds the request parameters for the UnlockDevice method.
type UnlockDeviceRequest struct {
	EventID    uuid.UUID `json:"event_id"`
	DeviceID   uuid.UUID `json:"device_id"`
	UnlockCode string    `json:"unlock_code"`
}

// UnlockDeviceResponse holds the response values for the UnlockDevice method.
type UnlockDeviceResponse struct {
	Session frontend.Session
	err     error
}

// Failed implements Failer.
func (r UnlockDeviceResponse) Failed() error { return r.err }

// GenerateQRRequest holds the request parameters for the GenerateQR method.
type GenerateQRRequest struct {
	EventID    uuid.UUID
	DeviceID   uuid.UUID
	UnlockCode string
}

// GenerateQRResponse holds the response values for the GenerateQR method.
type GenerateQRResponse struct {
	QR  []byte
	err error
}

// Failed implements Failer.
func (r GenerateQRResponse) Failed() error { return r.err }
