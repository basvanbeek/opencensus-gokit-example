package sd

import (
	// stdlib
	"io"
	"sort"
	"sync"
	"time"

	// external

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
)

// clientInstancerCache collects the most recent set of instances from a service
// discovery system, creates clientInstances for them using a factory function,
// and makes them available to consumers.
type clientInstancerCache struct {
	options            clientInstancerOptions
	mtx                sync.RWMutex
	factory            Factory
	cache              map[string]clientInstanceCloser
	err                error
	clientInstances    []interface{}
	logger             log.Logger
	invalidateDeadline time.Time
	timeNow            func() time.Time
}

type clientInstanceCloser struct {
	ci interface{}
	io.Closer
}

// newClientInstancerCache returns a new, empty clientInstancerCache.
func newClientInstancerCache(
	factory Factory, logger log.Logger, options clientInstancerOptions,
) *clientInstancerCache {
	return &clientInstancerCache{
		options: options,
		factory: factory,
		cache:   map[string]clientInstanceCloser{},
		logger:  logger,
		timeNow: time.Now,
	}
}

// Update should be invoked by clients with a complete set of current instance
// strings whenever that set changes. The cache manufactures new clientInstances
// via the factory, closes old clientInstances when they disappear, and persists
// existing clientInstances if they survive through an update.
func (c *clientInstancerCache) Update(event sd.Event) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// Happy path.
	if event.Err == nil {
		c.updateCache(event.Instances)
		c.err = nil
		return
	}

	// Sad path. Something's gone wrong in sd.
	c.logger.Log("err", event.Err)
	if !c.options.invalidateOnError {
		return // keep returning the last known clientInstances on error
	}
	if c.err != nil {
		return // already in the error state, do nothing & keep original error
	}
	c.err = event.Err
	// set new deadline to invalidate clientInstances unless non-error Event is
	// received
	c.invalidateDeadline = c.timeNow().Add(c.options.invalidateTimeout)
	return
}

func (c *clientInstancerCache) updateCache(instances []string) {
	// Deterministic order (for later).
	sort.Strings(instances)

	// Produce the current set of services.
	cache := make(map[string]clientInstanceCloser, len(instances))
	for _, instance := range instances {
		// If it already exists, just copy it over.
		if sc, ok := c.cache[instance]; ok {
			cache[instance] = sc
			delete(c.cache, instance)
			continue
		}

		// If it doesn't exist, create it.
		service, closer, err := c.factory(instance)
		if err != nil {
			c.logger.Log("instance", instance, "err", err)
			continue
		}
		cache[instance] = clientInstanceCloser{service, closer}
	}

	// Close any leftover clientInstances.
	for _, sc := range c.cache {
		if sc.Closer != nil {
			sc.Closer.Close()
		}
	}

	// Populate the slice of clientInstances.
	clientInstances := make([]interface{}, 0, len(cache))
	for _, instance := range instances {
		// A bad factory may mean an instance is not present.
		if _, ok := cache[instance]; !ok {
			continue
		}
		clientInstances = append(clientInstances, cache[instance].ci)
	}

	// Swap and trigger GC for old copies.
	c.clientInstances = clientInstances
	c.cache = cache
}

// Clients yields the current set of (presumably identical) clientInstances,
// ordered lexicographically by the corresponding instance string.
func (c *clientInstancerCache) Clients() ([]interface{}, error) {
	// in the steady state we're going to have many goroutines calling Clients()
	// concurrently, so to minimize contention we use a shared R-lock.
	c.mtx.RLock()

	if c.err == nil || c.timeNow().Before(c.invalidateDeadline) {
		defer c.mtx.RUnlock()
		return c.clientInstances, nil
	}

	c.mtx.RUnlock()

	// in case of an error, switch to an exclusive lock.
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// re-check condition due to a race between RUnlock() and Lock().
	if c.err == nil || c.timeNow().Before(c.invalidateDeadline) {
		return c.clientInstances, nil
	}

	c.updateCache(nil) // close any remaining active clientInstances
	return nil, c.err
}
