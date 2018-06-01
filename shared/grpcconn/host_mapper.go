package grpcconn

import (
	// stdlib
	"io"
	"sync"

	// external
	"google.golang.org/grpc"
)

// HostMapper manages a map of discovered instances with accompanying gRPC
// client connections.
type HostMapper interface {
	Get(instance string) (*grpc.ClientConn, io.Closer, error)
}

// NewHostMapper initializes and returns a HostMapper for a gRPC service.
func NewHostMapper(dialOptions ...grpc.DialOption) HostMapper {
	return &hostMapper{
		host:        make(map[string]*grpc.ClientConn),
		DialOptions: dialOptions,
	}
}

// hostMapper implements HostMapper
type hostMapper struct {
	mtx         sync.Mutex
	host        map[string]*grpc.ClientConn
	DialOptions []grpc.DialOption
}

// Get a gRPC client connection for provided instance.
func (h *hostMapper) Get(instance string) (*grpc.ClientConn, io.Closer, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if conn := h.host[instance]; conn != nil {
		return conn, &closer{h, instance}, nil
	}

	conn, err := grpc.Dial(instance, h.DialOptions...)
	if err != nil {
		return nil, nil, err
	}

	h.host[instance] = conn

	return conn, &closer{h, instance}, nil
}

// remove closes a gRPC client connection and removes it from the map.
func (h *hostMapper) remove(instance string) {
	h.mtx.Lock()
	if conn := h.host[instance]; conn != nil {
		defer conn.Close()
		delete(h.host, instance)
	}
	h.mtx.Unlock()
}

type closer struct {
	hm       *hostMapper
	instance string
}

// Close implements io.Closer
func (c *closer) Close() error {
	c.hm.remove(c.instance)
	return nil
}
