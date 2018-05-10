package routes

import (
	// external
	"github.com/gorilla/mux"
)

// Endpoints holds all available HTTP endpoints for our service.
type Endpoints struct {
	UnlockDevice *mux.Route
	GenerateQR   *mux.Route
}

// InitEndpoints wires the HTTP endpoints to our Go kit service endpoints.
func InitEndpoints(router *mux.Router) Endpoints {
	return Endpoints{
		UnlockDevice: router.
			Methods("POST").
			Path("/unlock_device/{event_id}/{device_id}").
			Queries("code", "{code}").
			Name("unlock_device"),
		GenerateQR: router.
			Methods("GET").
			Path("/generate_qr/{event_id}/{device_id}").
			Queries("code", "{code}").
			Name("generate_qr"),
	}
}
