package sd

import (
	// stdlib
	"time"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
)

// ClientInstancer listens to a service discovery system and yields a set of
// identical client instances on demand. An error indicates a problem with
// connectivity to the service discovery system, or within the system itself;
// an ClientInstancer may yield no client instances without error.
type ClientInstancer interface {
	Clients() ([]interface{}, error)
}

// FixedClientInstancer yields a fixed set of client instances.
type FixedClientInstancer []interface{}

// Clients implements ClientInstancer.
func (s FixedClientInstancer) Clients() ([]interface{}, error) { return s, nil }

// NewClientInstancer creates a ClientInstancer that subscribes to updates from
// Instancer src and uses factory f to create Client instances. If src notifies
// of an error, the ClientInstancer keeps returning previously created Client
// instances assuming they are still good, unless this behavior is disabled via
// InvalidateOnError option.
func NewClientInstancer(
	src sd.Instancer, f Factory, logger log.Logger, options ...Option,
) *DefaultClientInstancer {
	opts := clientInstancerOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	se := &DefaultClientInstancer{
		cache:     newClientInstancerCache(f, logger, opts),
		instancer: src,
		ch:        make(chan sd.Event),
	}
	go se.receive()
	src.Register(se.ch)
	return se
}

// Option allows control of clientCache behavior.
type Option func(*clientInstancerOptions)

// InvalidateOnError returns Option that controls how the ClientInstancer
// behaves when then Instancer publishes an Event containing an error.
// Without this option the ClientInstancer continues returning the last known
// client instances. With this option, the ClientInstancer continues returning
// the last known client instances until the timeout elapses, then closes all
// active client instances and starts returning an error. Once the Instancer
// sends a new update with valid resource instances, the normal operation is
// resumed.
func InvalidateOnError(timeout time.Duration) Option {
	return func(opts *clientInstancerOptions) {
		opts.invalidateOnError = true
		opts.invalidateTimeout = timeout
	}
}

type clientInstancerOptions struct {
	invalidateOnError bool
	invalidateTimeout time.Duration
}

// DefaultClientInstancer implements an ClientInstancer interface.
// When created with NewClientInstancer function, it automatically registers
// as a subscriber to events from the Instances and maintains a list
// of active client instances.
type DefaultClientInstancer struct {
	cache     *clientInstancerCache
	instancer sd.Instancer
	ch        chan sd.Event
}

func (dc *DefaultClientInstancer) receive() {
	for event := range dc.ch {
		dc.cache.Update(event)
	}
}

// Close deregisters DefaultClientInstancer from the Instancer and stops the
// internal go-routine.
func (dc *DefaultClientInstancer) Close() {
	dc.instancer.Deregister(dc.ch)
	close(dc.ch)
}

// Clients implements ClientInstancer.
func (dc *DefaultClientInstancer) Clients() ([]interface{}, error) {
	return dc.cache.Clients()
}
