package sd

import (
	// stdlib
	"sync/atomic"
)

// NewRoundRobin returns a load balancer that returns services in sequence.
func NewRoundRobin(s ClientInstancer) Balancer {
	return &roundRobin{
		s: s,
		c: 0,
	}
}

type roundRobin struct {
	s ClientInstancer
	c uint64
}

func (rr *roundRobin) Client() (interface{}, error) {
	clients, err := rr.s.Clients()
	if err != nil {
		return nil, err
	}
	if len(clients) <= 0 {
		return nil, ErrNoClients
	}
	old := atomic.AddUint64(&rr.c, 1) - 1
	idx := old % uint64(len(clients))
	return clients[idx], nil
}
