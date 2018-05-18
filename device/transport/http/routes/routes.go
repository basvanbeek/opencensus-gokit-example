package routes

import (
	// external
	"github.com/gorilla/mux"
)

// Endpoints holds all available HTTP endpoints for our service.
type Endpoints struct {
	Unlock *mux.Route
}

// InitEndpoints wires the HTTP endpoints to our Go kit service endpoints.
func InitEndpoints(router *mux.Router) Endpoints {
	return Endpoints{
		Unlock: router.
			Methods("POST").
			Path("/unlock/{event_id}/{device_id}").
			Queries("code", "{code}").
			Name("unlock"),
	}
}
