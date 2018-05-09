package transport

import "github.com/basvanbeek/opencensus-gokit-example/qr"

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// GenerateRequest holds the request parameters for the Generate method
type GenerateRequest struct {
	Data  string
	Level qr.RecoveryLevel
	Size  int
}

// GenerateResponse holds the response values for the Generate method
type GenerateResponse struct {
	QR  []byte
	err error
}

// Failed implements Failer
func (r GenerateResponse) Failed() error { return r.err }
