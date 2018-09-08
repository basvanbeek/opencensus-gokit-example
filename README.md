# opencensus-gokit-example
Example of using OpenCensus with Go kit

[![Go Report Card](https://goreportcard.com/badge/github.com/basvanbeek/opencensus-gokit-example)](https://goreportcard.com/report/github.com/basvanbeek/opencensus-gokit-example)

# notable dependencies
- Go kit
- OpenCensus
- gRPC
- Twirp
- Zipkin

# build services

Services can be built by running go generate in the root of the project:

```sh
$ go generate ./...
```

# running the services

Get a Zipkin instance running by following the directions as found
[here](https://zipkin.io/pages/quickstart). This will allow the services
instrumented with OpenCensus to export tracing details to a Zipkin backend.

Get etcd running for service discovery. Instructions to get it up and
running can be found [here](https://coreos.com/etcd/docs/latest/dl_build.html).

Now you can start the various services included in this demo.

Example:
```sh
#!/bin/sh
nohup ./ocg-qrgenerator &>qrgenerator.log &
nohup ./ocg-device      &>device.log      &
nohup ./ocg-event       &>event.log       &
nohup ./ocg-frontend    &>frontend.log    &
```

Each service will dynamically select available ports to listen on and advertise
these on etcd. It is possible to run multiple instances for each service on a
single machine. The clients can automatically load balance and retry on the
available services.
