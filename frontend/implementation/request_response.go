package implementation

import (
	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/satori/go.uuid"
)

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

type UnlockDeviceRequest struct {
	EventID    uuid.UUID `json:"event_id"`
	DeviceID   uuid.UUID `json:"device_id"`
	UnlockCode string    `json:"unlock_code"`
}

type UnlockDeviceResponse struct {
	Session frontend.Session
	err     error
}

func (r UnlockDeviceResponse) Failed() error { return r.err }

type GenerateQRRequest struct {
	EventID    uuid.UUID
	DeviceID   uuid.UUID
	UnlockCode string
}

type GenerateQRResponse struct {
	QR  []byte
	err error
}

func (r GenerateQRResponse) Failed() error { return r.err }
