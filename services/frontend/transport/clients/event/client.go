package event

import (
	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	eventtwirp "github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/clients/event/twirp"
)

// NewTwirp returns a new event client using the Twirp transport.
func NewTwirp(instancer sd.Instancer, logger log.Logger) event.Service {
	return eventtwirp.NewClient(instancer, logger)
}
