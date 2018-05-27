package sd

import (
	"errors"
)

// Balancer yields client instances according to some heuristic.
type Balancer interface {
	Client() (interface{}, error)
}

// ErrNoClients is returned when no qualifying client instances are available.
var ErrNoClients = errors.New("no client instance available")
