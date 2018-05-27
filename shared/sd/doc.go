// Package sd provides utilities related to service discovery. That includes the
// client-side loadbalancer pattern, where a microservice subscribes to a
// service discovery system in order to reach remote instances; as well as the
// registrator pattern, where a microservice registers itself in a service
// discovery system. Implementations are provided for most common systems.
//
// Most of the code in this package is taken from the awesome Go kit library and
// adjusted so generic clients can be created and no dependency on the Go kit
// Endpoint concept is needed. Where possible it uses the Go kit sd package.
//
// See https:/gokit.io for more.
package sd
