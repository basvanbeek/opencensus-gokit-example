package sd

import "io"

// Factory is a function that converts an instance string (e.g. host:port) to a
// specific client implementation. A factory also returns an io.Closer that's
// invoked when the instance goes away and needs to be cleaned up. Factories may
// return nil closers.
//
// Users are expected to provide their own factory functions that assume
// specific transports, or can deduce transports by parsing the instance string.
type Factory func(instance string) (interface{}, io.Closer, error)
