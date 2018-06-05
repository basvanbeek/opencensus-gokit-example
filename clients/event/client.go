package event

import (
	// external
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"

	// project
	eventtwirp "github.com/basvanbeek/opencensus-gokit-example/clients/event/twirp"
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
)

// NewTwirp returns a new event client using the Twirp transport.
func NewTwirp(instancer sd.Instancer, client *http.Client, logger log.Logger) event.Service {
	return eventtwirp.NewClient(instancer, client, logger)
}
