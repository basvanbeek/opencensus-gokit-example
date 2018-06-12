package transport

import (
	// external
	"github.com/go-kit/kit/endpoint"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
)

var (
	_ endpoint.Failer = GenerateResponse{}
)

// GenerateRequest holds the request parameters for the Generate method.
type GenerateRequest struct {
	Data  string
	Level qr.RecoveryLevel
	Size  int
}

// GenerateResponse holds the response values for the Generate method.
type GenerateResponse struct {
	QR  []byte
	Err error
}

// Failed implements Failer.
func (r GenerateResponse) Failed() error { return r.Err }
