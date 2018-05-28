package routes

import (
	// external
	"github.com/gorilla/mux"
)

// Endpoints holds all available HTTP endpoints for our service.
type Endpoints struct {
	Login        *mux.Route
	EventCreate  *mux.Route
	EventGet     *mux.Route
	EventUpdate  *mux.Route
	EventDelete  *mux.Route
	EventList    *mux.Route
	UnlockDevice *mux.Route
	GenerateQR   *mux.Route
}

// InitEndpoints wires the HTTP endpoints to our Go kit service endpoints.
func InitEndpoints(router *mux.Router) Endpoints {
	return Endpoints{
		Login: router.
			Methods("POST").
			Path("/login").
			Name("login"),
		EventCreate: router.
			Methods("POST").
			Path("/event").
			Name("event_create"),
		EventGet: router.
			Methods("GET").
			Path("/event/{event_id}").
			Name("event_get"),
		EventUpdate: router.
			Methods("PUT").
			Path("/event/{event_id}").
			Name("event_update"),
		EventDelete: router.
			Methods("DELETE").
			Path("/event/{event_id}").
			Name("event_delete"),
		EventList: router.
			Methods("GET").
			Path("/event").
			Name("event_list"),
		UnlockDevice: router.
			Methods("POST", "GET").
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
