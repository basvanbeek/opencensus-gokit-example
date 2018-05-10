package transport

import (
	// project
	"github.com/basvanbeek/opencensus-gokit-example/qr"
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
